package auth

import (
	xadmin "github.com/wgbbiao/goxadmin"

	validator "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func initValidator() {
	validate = validator.New()
	validate.RegisterStructValidation(CreateUserStructLevelValidation, &User{})
}

//CreateUserStructLevelValidation CreateUserStructLevelValidation
func CreateUserStructLevelValidation(sl validator.StructLevel) {
	j := sl.Current().Interface().(User)
	if j.Password != j.Password2 {
		sl.ReportError(j.Password, "Password", "Password", xadmin.UserPasswordError, "")
		return
	}
	j.Password = "22222"
}
