package main

import (
	"fmt"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql
	"github.com/kataras/iris/v12"
	"github.com/wgbbiao/goxadmin"
)

func main() {
	DB, _ := gorm.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local",
		"root",
		"123456",
		"tcp(192.168.3.158:3306)",
		"app_rowclub"))

	DB.LogMode(true)

	DB.SingularTable(true)

	DB.DB().Ping()
	DB.DB().SetMaxIdleConns(50)
	DB.DB().SetMaxOpenConns(50)
	DB.DB().SetConnMaxLifetime(time.Duration(1000) * time.Second)
	// time.LoadLocation(cfg.Section("system").Key("location").String())
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "Origin"},
	})

	r := iris.New()
	r.Use(crs)
	r.Use(func(ctx iris.Context) {
		ctx.Gzip(true)
		ctx.Next()
	})
	r.Options("{root:path}", func(context iris.Context) {
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Header("Access-Control-Allow-Headers", "Origin,Authorization,Content-Type,Accept,X-Total,X-Limit,X-Offset")
		context.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,HEAD")
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Expose-Headers", "Content-Length,Content-Encoding,Content-Type")
	})
	xadmin := goxadmin.NewXadmin(DB, r.Party("/admin"))
	xadmin.Init()
	for _, _r := range r.GetRoutes() {
		fmt.Println(_r)
	}
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
