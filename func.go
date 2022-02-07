package goxadmin

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

//Cmd5 Cmd5
func Cmd5(txt, salt string) string {
	m5 := md5.New()
	m5.Write([]byte(txt))
	m5.Write([]byte(salt))
	st := m5.Sum(nil)
	return hex.EncodeToString(st)
}

//SyncPermissions 同步权限,同步
func (o *XadminConfig) SyncPermissions() {
	acts := []string{
		PolicyCreate,
		PolicyView,
		PolicyDelete,
		PolicyUpdate,
	}
	for _, model := range o.GetRegModels() {
		v := reflect.ValueOf(model)
		method := v.MethodByName("Permissions")
		newActs := acts
		if method.Kind() == reflect.Func {
			args := make([]reflect.Value, 0)
			values := method.Call(args)
			dd := values[0].Interface()
			newActs = append(newActs, dd.([]string)...)
		}
		table := GetModelName(model)
		for _, act := range newActs {
			AddPermission(table, act)
		}
	}
}

//GetVerboseName 取得models名称
func GetVerboseName(m interface{}) string {
	v := reflect.ValueOf(m)
	method := v.MethodByName("VerboseName")
	if method.Kind() == reflect.Func {
		args := make([]reflect.Value, 0)
		values := method.Call(args)
		return values[0].String()
	}
	return v.Elem().Type().Name()
}

//GetModelName 取得model名称
func GetModelName(m interface{}) (ct ContentType) {
	path := reflect.TypeOf(m).String()
	path = strings.Replace(path, "*", "", 1)
	paths := strings.Split(path, ".")
	ct.AppLabel = paths[0]
	ct.Model = paths[1]
	Db.FirstOrCreate(&ct, ct)
	return
}

//GetActionByMethod 取得权限的名称
func GetActionByMethod(method string) (action string) {
	switch method {
	case iris.MethodGet:
		action = PolicyView
	case iris.MethodPost:
		action = PolicyCreate
	case iris.MethodPut:
		action = PolicyUpdate
	case iris.MethodDelete:
		action = PolicyDelete
	default:
		action = ""
	}
	return
}

//GenCodeName 取得权限名称
func GenCodeName(code, modelname string) string {
	return fmt.Sprintf("%s_%s", code, strings.ToLower(modelname))
}

//GetConfig 取得配置文件
func GetConfig(model, table string) Config {
	return Xadmin.Models[fmt.Sprintf("%s/%s", model, table)]
}

//ShowJSON 格式化json
func ShowJSON(data interface{}) {
	b, _ := json.Marshal(data)
	var str bytes.Buffer
	_ = json.Indent(&str, b, "", "    ")
	fmt.Println("formated: ", str.String())
}

//MapToWhere map转成gorm需要的搜索条件
func MapToWhere(params map[string]string, config Config) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range params {
			if strings.HasPrefix(k, "_p_") {
				k = strings.Replace(k, "_p_", "", -1)
				if strings.Contains(k, ".") {
					tmp := strings.Split(k, ".")
					tableName := tmp[0]
					db = db.Joins(tableName)
				}
				fields := strings.Split(k, "__")
				field := fields[0]
				paramType := fields[1]
				switch paramType {
				case "exact":
					db = db.Where(fmt.Sprintf("%s = ?", field), v)
				case "in":
					db = db.Where(fmt.Sprintf("%s in (?)", field), strings.Split(v, ","))
				case "notin":
					db = db.Where(fmt.Sprintf("%s not in (?)", field), strings.Split(v, ","))
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
		order, ok := params["o"]
		if !ok {
			if config.Sort != "" {
				order = config.Sort
			} else {
				order = "-id"
			}
		}
		if strings.HasPrefix(order, "-") {
			db = db.Order(fmt.Sprintf("%s DESC", order[1:]))
		} else {
			db = db.Order(fmt.Sprintf("%s ASC", order))
		}
		preloads, ok := params["preloads"]
		if !ok {
			db = db.Set("gorm:auto_preload", false)
		} else {
			for _, p := range strings.Split(preloads, ",") {
				db = db.Preload(p)
			}
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
