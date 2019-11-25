package auth

import (
	"reflect"
	xadmin "github.com/wgbbiao/goxadmin"
	"strings"
)

//常量
const (
	//权限
	PolicyRead   string = "view"
	PolicyWrite  string = "edit"
	PolicyCreate string = "create"
	PolicyDelete string = "delete"
)

//SyncPermissions 同步权限
func SyncPermissions() {
	acts := []string{
		PolicyCreate,
		PolicyRead,
		PolicyDelete,
		PolicyWrite,
	}
	for _, model := range xadmin.GetRegModels() {
		v := reflect.ValueOf(model)
		method := v.MethodByName("Permissions")
		newActs := acts
		if method.Kind() == reflect.Func {
			args := make([]reflect.Value, 0)
			values := method.Call(args)
			// 	for _, val := range values[0].([]) {
			// 		fmt.Println(val.String())
			// 		newActs = append(newActs, val.String())
			// 	}
			// fmt.Println(values[0])
			dd := values[0].Interface()
			newActs = append(newActs, dd.([]string)...)
		}
		modelname := GetVerboseName(model)
		table := GetModelName(model)
		for _, act := range newActs {
			AddPermission(table, act, modelname)
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
func GetModelName(m interface{}) string {
	path := reflect.TypeOf(m).String()
	path = strings.Replace(path, "*", "", 1)
	return path
}
