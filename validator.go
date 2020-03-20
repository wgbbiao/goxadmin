package goxadmin

import (
	validator "gopkg.in/go-playground/validator.v9"
)

//Validate 表单验证
var Validate *validator.Validate

func initValidator() {
	Validate = validator.New()
	Validate.RegisterStructValidation(CreateUserStructLevelValidation, &User{})
}

//RegisterStructValidation RegisterStructValidation
func RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{}) {
	Validate.RegisterStructValidation(fn, types...)
}

//CreateUserStructLevelValidation CreateUserStructLevelValidation
func CreateUserStructLevelValidation(sl validator.StructLevel) {
	j := sl.Current().Interface().(User)
	if j.Password != j.Password2 {
		sl.ReportError(j.Password, "Password", "Password", UserPasswordError, "")
		return
	}
}
