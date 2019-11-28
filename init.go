package goxadmin

import (
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

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

//Db Db
var Db *gorm.DB

var models map[string]Config

//Model 默认Model
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
