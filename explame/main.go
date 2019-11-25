package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql
	"github.com/kataras/iris/v12"
	"github.com/wgbbiao/goxadmin"
	"github.com/wgbbiao/goxadmin/auth"
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
	goxadmin.Db = DB
	auth.AutoMigrate()     //生成表结构
	auth.SyncPermissions() //同步权限
	r := iris.New()
	goxadmin.Init(r)

	for _, _r := range r.GetRoutes() {
		fmt.Println(_r)
	}
	r.Run(iris.Addr("127.0.0.1:8099"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              true,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         true,
		TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
		Charset:                           "UTF-8",
	}))
}
