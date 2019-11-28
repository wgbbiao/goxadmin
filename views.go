package goxadmin

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

//Login 用户登录
func Login(c iris.Context) {
	type Form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var form Form
	var u User
	if err := c.ReadJSON(&form); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{
			"status": HTTPFail,
			"code":   FormReadError,
		})
	} else {
		if db := u.GetByUsername(form.Username); db.RecordNotFound() {
			c.StatusCode(iris.StatusBadRequest)
			c.JSON(iris.Map{
				"status": HTTPFail,
				"code":   UserDoesNotExist,
			})
		} else {
			if u.CheckPassword(form.Password) {
				claim := jwt.MapClaims{
					"exp": time.Now().Unix() + 86400,
					"uid": u.ID,
				}

				accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
				tokenString, _ := accessToken.SignedString([]byte(JwtKey))
				u.UpdateInfo(map[string]interface{}{
					"last_login_at": time.Now(),
				})
				c.JSON(iris.Map{
					"status":   HTTPSuccess,
					"token":    tokenString,
					"username": u.Username,
				})
			} else {
				c.StatusCode(iris.StatusBadRequest)
				c.JSON(iris.Map{
					"status": HTTPFail,
					"code":   UserPasswordError,
				})
			}
		}
	}
}

//GetInfo 取得用户信息
func GetInfo(c iris.Context) {
	u := c.Values().Get("u").(User)
	models, permissions := u.GetPermissionInfo()
	c.JSON(iris.Map{
		"status":      "success",
		"username":    u.Username,
		"isSuper":     u.IsSuper,
		"models":      models,
		"permissions": permissions,
	})
}

var jwtc = jwtmiddleware.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler:  OnJwtError,
}

var myJwtMiddleware = jwtmiddleware.New(jwtc)

//OnJwtError jwt error
func OnJwtError(ctx iris.Context, err error) {
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(iris.Map{
		"status": "fail",
		"info":   err,
		"code":   TokenIsExpired,
	})
}

//CheckJWTAndSetUser 检查jwt并把User放到Values
func CheckJWTAndSetUser(ctx iris.Context) {
	if err := myJwtMiddleware.CheckJWT(ctx); err != nil {
		myJwtMiddleware.Config.ErrorHandler(ctx, err)
		return
	}
	// If everything ok then call next.
	if ctx.GetStatusCode() != iris.StatusUnauthorized {
		var u *User
		x, _ := ctx.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
		if rt := u.GetUserByID(int(x["uid"].(float64))); !rt.RecordNotFound() && rt.Error == nil {
			config := GetConfig(ctx.Params().Get("model"), ctx.Params().GetString("table"))
			bl := HasPermissionForModel(u, config.Model, GetActionByMethod(ctx.Method()))
			if bl {
				ctx.Values().Set("u", u)
				ctx.Next()
			} else {
				ctx.StatusCode(iris.StatusForbidden)
				ctx.JSON(iris.Map{
					"status": HTTPForbidden,
					"code":   UserNoPermission,
				})
			}
		} else {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status": HTTPFail,
				"code":   UserDoesNotExist,
			})
		}
	}
}
