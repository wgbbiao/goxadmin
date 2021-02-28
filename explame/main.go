package main

import (
	"fmt"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/wgbbiao/goxadmin"
	"gorm.io/driver/mysql" //mysql
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DB db
var DB *gorm.DB

// User 有多张 CreditCard，UserID 是外键
type User struct {
	gorm.Model
	Name        string
	CreditCards []CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}

var xadmin *goxadmin.XadminConfig

func init() {

	dsn := "root:123456@tcp(192.168.1.7:3306)/app_rowclub?charset=utf8mb4&parseTime=True&loc=Local"
	DB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	// DB.LogMode(true)

	sqldb, _ := DB.DB()
	sqldb.SetMaxIdleConns(50)
	sqldb.SetMaxOpenConns(50)
	sqldb.SetConnMaxLifetime(time.Duration(1000) * time.Second)

	xadmin = goxadmin.NewXadmin(DB)
	xadmin.Register(&User{}, goxadmin.Config{})
	xadmin.Register(&CreditCard{}, goxadmin.Config{})

	xadmin.AutoMigrate()
	xadmin.SyncPermissions()
}
func main() {
	// fmt.Println("asdfasdf")
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
		// ctx.Gzip(true)
		ctx.Next()
	})
	r.Options("{root:path}", func(context iris.Context) {
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Header("Access-Control-Allow-Headers", "Origin,Authorization,Content-Type,Accept,X-Total,X-Limit,X-Offset")
		context.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,HEAD")
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Expose-Headers", "Content-Length,Content-Encoding,Content-Type")
	})
	xadmin.SetIris(r.Party("/admin"))
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
}
