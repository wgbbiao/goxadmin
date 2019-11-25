package auth

import (
	xadmin "github.com/wgbbiao/goxadmin"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
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
			"status": xadmin.HTTPFail,
			"code":   xadmin.FormReadError,
		})
	} else {
		if db := u.GetByUsername(form.Username); db.RecordNotFound() {
			c.StatusCode(iris.StatusBadRequest)
			c.JSON(iris.Map{
				"status": xadmin.HTTPFail,
				"code":   xadmin.UserDoesNotExist,
			})
		} else {
			if u.CheckPassword(form.Password) {
				claim := jwt.MapClaims{
					"exp": time.Now().Unix() + 86400,
					"uid": u.ID,
				}

				accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
				tokenString, _ := accessToken.SignedString([]byte(xadmin.JwtKey))
				u.UpdateInfo(map[string]interface{}{
					"last_login_at": time.Now(),
				})
				c.JSON(iris.Map{
					"status":   xadmin.HTTPSuccess,
					"token":    tokenString,
					"username": u.Username,
				})
			} else {
				c.StatusCode(iris.StatusBadRequest)
				c.JSON(iris.Map{
					"status": xadmin.HTTPFail,
					"code":   xadmin.UserPasswordError,
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
