package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql
	"github.com/kataras/iris/v12"
	"github.com/wgbbiao/goxadmin"
)

func main() {
	DB, _ := gorm.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local",
		"root",
		"123456",
		"tcp(192.168.1.153:3306)",
		"app_rowclub"))

	DB.LogMode(true)

	DB.SingularTable(true)

	DB.DB().Ping()
	DB.DB().SetMaxIdleConns(50)
	DB.DB().SetMaxOpenConns(50)
	DB.DB().SetConnMaxLifetime(time.Duration(1000) * time.Second)
	r := iris.New()
	xadmin := goxadmin.NewXadmin(DB, r, goxadmin.CheckJWTAndSetUser)
	xadmin.Init()
	r.Run(iris.Addr(":9999"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         true,
		TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
		Charset:                           "UTF-8",
	}))
	// dd := DB.NewScope(goxadmin.User{})
	// field, ok := dd.FieldByName("Permissions")
	// fmt.Println(ok)
	// b, _ := json.Marshal(field.Relationship)
	// var str bytes.Buffer
	// _ = json.Indent(&str, b, "", "    ")
	// fmt.Println("formated: ", str.String())
	// fmt.Println("data: ", string(b))

	// kids := make([]goxadmin.User, 0)
	// params := make(map[string]string)
	// params["_p_permissions.name__like"] = "view"
	// DB.Scopes(goxadmin.MapToWhere(params, goxadmin.Config{
	// 	Model: &goxadmin.User{},
	// })).Find(&kids)
	// // models.DB.Joins("left join user on kid.user_id = user.id").
	// // Where("user.mobile = ?", "13466625910").Find(&kids)
	// fmt.Println(kids)
}
