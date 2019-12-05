package goxadmin

//常量
const (
	HTTPSuccess   string = "succ"
	HTTPFail      string = "fail"
	HTTPForbidden string = "Forbidden"
)

//错误常量
const (
	DBError           string = "db_error"
	ValidateError     string = "validate_error"
	FormReadError     string = "form_read_error"
	TokenIsExpired    string = "token_is_expired"
	UserDoesNotExist  string = "user_does_not_exist"
	UserPasswordError string = "user_password_error"
	UserNoPermission  string = "user_no_permission"
)

//权限
const (
	PolicyView   string = "view"
	PolicyUpdate string = "update"
	PolicyCreate string = "create"
	PolicyDelete string = "delete"
)
