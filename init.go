package goxadmin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
	"github.com/wxnacy/wgo/arrays"

	"github.com/kataras/iris/v12"
)

//JwtKey JwtKey
var JwtKey string = "Auys7;fq272/csH6"

// JwtCheckFunc JwtCheckFunc
var JwtCheckFunc func(c iris.Context)

//Handle Handle
type Handle struct {
	Path   string
	Method []string
	Func   func(ctx iris.Context)
	Jwt    bool
}

//Handles 自定义handle
var Handles []Handle

// Config Config
type Config struct {
	ListField []string //列表页字段
	Model     interface{}
	// Layout     interface{} //表单排列
	PageSize   int //每页大小
	BeforeSave func(obj interface{})
	//todo
	// Form interface{} //自定表单
	Sort          string
	DisableAction []string //禁止的操作 [create,update,detail,delete,list]
}

//HasPremisson 检查是否有权限
func (c *Config) HasPremisson(action string) bool {
	return true
}

//Title Title
func (c *Config) Title() string {
	return "Title"
}

//Site 网站配置信息
type Site struct {
	Title  string //网站标题
	models map[string]Config
}

//GetMenu 网站标题 循环models
func (s *Site) GetMenu() {

}

//Db Db
var Db *gorm.DB

var models map[string]Config

//Model 默认Model
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//GetConfig 取得配置文件
func GetConfig(model, table string) Config {
	return models[fmt.Sprintf("%s/%s", model, table)]
}

// Menu 取得菜单
func (s *Site) Menu() interface{} {
	menus := make([]interface{}, 0)
	for model, conf := range models {
		menus = append(menus, map[string]interface{}{
			"title": conf.Title(),
			"model": model,
		})
	}
	return menus
}

//MapToWhere map转成gorm需要的搜索条件
func MapToWhere(status map[string]string, config Config) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range status {
			if strings.HasPrefix(k, "_p_") {
				k = strings.Replace(k, "_p_", "", -1)
				fields := strings.Split(k, "__")
				field := fields[0]
				paramType := fields[1]
				switch paramType {
				case "exact":
					db = db.Where(fmt.Sprintf("%s = ?", field), v)
				case "in":
					db = db.Where(fmt.Sprintf("%s in (?)", field), strings.Split(v, ","))
				case "to":
					db = db.Where(fmt.Sprintf("%s <= ?", field), v)
				case "from":
					db = db.Where(fmt.Sprintf("%s >= ?", field), v)
				case "not":
					db = db.Not(field, v)
				case "null":
					if v == "true" {
						db = db.Where(fmt.Sprintf("%s is null", field))
					} else {
						db = db.Where(fmt.Sprintf("%s is not null", field))
					}
				case "like":
					db = db.Where(fmt.Sprintf("%s like ?", field), fmt.Sprintf("%%%s%%", v))
				}
			}
		}
		order, ok := status["o"]
		if ok == false {
			if config.Sort != "" {
				order = config.Sort
			} else {
				order = "id"
			}
		}
		if strings.HasPrefix(order, "-") {
			db = db.Order(gorm.Expr("? DESC", order[1:]))
		} else {
			db = db.Order(gorm.Expr("? ASC", order))
		}
		return db
	}
}

//NewSlice 新的数组
func NewSlice(model interface{}) interface{} {
	if model == nil {
		return nil
	}
	sliceType := reflect.SliceOf(reflect.TypeOf(model))
	slice := reflect.MakeSlice(sliceType, 0, 0)
	slicePtr := reflect.New(sliceType)
	slicePtr.Elem().Set(slice)
	return slicePtr.Interface()
}

//GetVal 取得 结构体的实例
func GetVal(model interface{}) interface{} {
	return reflect.New(reflect.TypeOf(model).Elem()).Interface()
}

// ListHandel ListHandel
func ListHandel(ctx iris.Context) {
	config := GetConfig(ctx.Params().Get("model"), ctx.Params().Get("table"))
	if arrays.ContainsString(config.DisableAction, "list") > -1 {
		ctx.StatusCode(iris.StatusForbidden)
	} else {
		rs := NewSlice(config.Model)
		page := com.StrTo(ctx.URLParam("page")).MustInt()
		all, _ := ctx.URLParamBool("__all__")
		if page == 0 {
			page = 1
		}
		limit := config.PageSize
		if all {
			limit = 999999
		} else {
			if limit == 0 {
				limit = 20
			}
		}
		offset := (page - 1) * limit
		params := ctx.URLParams()
		cnt := 0

		err := Db.Set("gorm:auto_preload", true).Scopes(MapToWhere(params, config)).
			Limit(limit).
			Offset(offset).
			Find(rs).
			Offset(0).
			Count(&cnt).Error
		if err == nil {
			ctx.JSON(iris.Map{
				"list":  rs,
				"total": cnt,
			})
		} else {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status":  HTTPFail,
				"error":   "save error",
				"errinfo": err,
			})
		}
	}
}

//DetailHandel 详情页
func DetailHandel(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")
	config := GetConfig(ctx.Params().Get("model"), ctx.Params().Get("table"))
	if arrays.ContainsString(config.DisableAction, "detail") > -1 {
		ctx.StatusCode(iris.StatusForbidden)
	} else {
		obj := GetVal(config.Model)
		if err := Db.Set("gorm:auto_preload", true).First(obj, id).Error; err == nil {
			ctx.JSON(iris.Map{
				"data": obj,
			})
		} else {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status":  HTTPFail,
				"error":   "save error",
				"errinfo": err,
			})
		}
	}
}

//PostHandel 添加记录
func PostHandel(ctx iris.Context) {
	config := GetConfig(ctx.Params().Get("model"), ctx.Params().Get("table"))
	if arrays.ContainsString(config.DisableAction, "create") > -1 {
		ctx.StatusCode(iris.StatusForbidden)
	} else {
		obj := GetVal(config.Model)
		if err := ctx.ReadJSON(&obj); err == nil {
			if err = Validate.Struct(obj); err == nil {
				if config.BeforeSave != nil {
					config.BeforeSave(obj)
				}
				if err = Db.Create(obj).Error; err == nil {
					ctx.JSON(iris.Map{
						"status": HTTPSuccess,
						"data":   obj,
					})
				} else {
					ctx.StatusCode(iris.StatusBadRequest)
					ctx.JSON(iris.Map{
						"status":  HTTPFail,
						"error":   DBError,
						"errinfo": err,
					})
				}
			} else {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"status":  HTTPFail,
					"error":   ValidateError,
					"errinfo": err,
				})
			}
		} else {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status": HTTPFail,
				"error":  FormReadError,
			})
		}
	}
}

//UpdateHandel 修改记录
func UpdateHandel(ctx iris.Context) {
	config := GetConfig(ctx.Params().Get("model"), ctx.Params().Get("table"))
	if arrays.ContainsString(config.DisableAction, "update") > -1 {
		ctx.StatusCode(iris.StatusForbidden)
	} else {
		obj := GetVal(config.Model)
		id, _ := ctx.Params().GetInt("id")
		Db.First(obj, id)
		t := reflect.TypeOf(obj).Elem()
		// 删除多对多的关系，然后重新添加
		for k := 0; k < t.NumField(); k++ {
			field := t.Field(k)
			if strings.Contains(field.Tag.Get("gorm"), "many2many") {
				Db.Model(obj).Association(t.Field(k).Name).Clear()
			}
		}

		if err := ctx.ReadJSON(&obj); err == nil {
			// config.BeforeSave(obj)
			if Db.Save(obj).Error == nil {
				sc := Db.NewScope(obj)
				for _, f := range sc.Fields() {
					if f.Relationship != nil {
						if f.Relationship.Kind == "has_many" {
							ff := f.Field
							px := make([]int64, 0)
							for index := 0; index < ff.Len(); index++ {
								elm := ff.Index(index)
								px = append(px, elm.FieldByName("ID").Int())
							}
							if ff.Len() > 0 {
								tableName := Db.NewScope(ff.Index(0).Interface()).TableName()
								sql := fmt.Sprintf("DELETE from %s Where id NOT IN (?) and %s = ?", tableName, f.Relationship.ForeignDBNames[0])
								Db.Exec(sql, px, id)
							}
						}
					}
				}
				ctx.JSON(iris.Map{
					"status": HTTPSuccess,
				})
			} else {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"status":  HTTPFail,
					"error":   DBError,
					"errinfo": err,
				})
			}
		} else {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status": HTTPFail,
				"error":  FormReadError,
			})
		}
	}
}

//DeleteHandel 删除记录
func DeleteHandel(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")
	config := GetConfig(ctx.Params().Get("model"), ctx.Params().Get("table"))
	if arrays.ContainsString(config.DisableAction, "delete") > -1 {
		ctx.StatusCode(iris.StatusForbidden)
	} else {
		obj := GetVal(config.Model)
		if err := Db.First(obj, id).Error; err == nil {
			if Db.Delete(obj).Error == nil {
				ctx.JSON(iris.Map{
					"data": obj,
				})
			} else {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(iris.Map{
					"status": HTTPFail,
					"error":  DBError,
				})
			}
		} else {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"errinfo": err,
			})
		}
	}
}

// //GetInfo 取得表结构（字段信息用于前台表单创建）
// func GetInfo(ctx iris.Context) {

// }

//Register Register
func Register(model interface{}, config Config) {
	if models == nil {
		models = make(map[string]Config)
	}
	resType := reflect.TypeOf(model)
	config.Model = model

	modelname := strings.Replace(resType.String(), "*", "", 1)
	modelname = strings.Replace(modelname, ".", "/", 1)
	models[modelname] = config
}

//GetRegModels GetRegModels
func GetRegModels() (ms []interface{}) {
	for _, conf := range models {
		ms = append(ms, conf.Model)
	}
	return
}

//Setdb 设置数据库链接
func Setdb(d *gorm.DB) {
	Db = d
}

//XadminIrisParty XadminIrisParty
var XadminIrisParty iris.Party

//RegisterView 注册新的api
func RegisterView(handle ...Handle) {
	Handles = append(Handles, handle...)
}

//Init Init
func Init(r iris.Party) {
	JwtCheckFunc = CheckJWTAndSetUser
	XadminIrisParty = r
	for _, handel := range Handles {
		for _, method := range handel.Method {
			if handel.Jwt {
				XadminIrisParty.Handle(method, handel.Path, JwtCheckFunc, handel.Func)
			} else {
				XadminIrisParty.Handle(method, handel.Path, handel.Func)
			}
		}
	}
	XadminIrisParty.Get("{model:string}/{table:string}", JwtCheckFunc, ListHandel)
	XadminIrisParty.Get("{model:string}/{table:string}/{id:int}", JwtCheckFunc, DetailHandel)
	XadminIrisParty.Put("{model:string}/{table:string}/{id:int}", JwtCheckFunc, UpdateHandel)
	XadminIrisParty.Post("{model:string}/{table:string}", JwtCheckFunc, PostHandel)
	XadminIrisParty.Delete("{model:string}/{table:string}/{id:int}", JwtCheckFunc, DeleteHandel)
}

//Cmd5 Cmd5
func Cmd5(txt, salt string) string {
	m5 := md5.New()
	m5.Write([]byte(txt))
	m5.Write([]byte(salt))
	st := m5.Sum(nil)
	return hex.EncodeToString(st)
}
