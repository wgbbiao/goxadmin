module xadmin

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/iris-contrib/middleware/jwt v0.0.0-20191111233249-6e1f5f92190e
	github.com/jinzhu/gorm v1.9.11
	github.com/kataras/iris v11.1.1+incompatible
	github.com/kataras/iris/v12 v12.0.1
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/unknwon/com v1.0.1
	github.com/wgbbiao/goxadmin v0.0.0-00010101000000-000000000000
	github.com/wxnacy/wgo v1.0.4
	gopkg.in/go-playground/validator.v9 v9.30.0
)

replace github.com/wgbbiao/goxadmin => ./
