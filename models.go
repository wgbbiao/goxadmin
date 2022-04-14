package goxadmin

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

//DefaultModel 默认Model
type DefaultModel struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//User 管理员
type User struct {
	DefaultModel
	Username    string        `gorm:"type:varchar(50);UNIQUE" json:"username"`
	Password    string        `gorm:"type:varchar(50)" json:"password,omitempty"`
	Password2   string        `gorm:"-" json:"password2,omitempty"`
	Salt        string        `gorm:"type:varchar(64)" json:"-"`
	IsSuper     bool          `gorm:"default:false" json:"is_super"`
	LastLoginAt *time.Time    `gorm:"type:datetime;null" json:"last_login_at"`
	Roles       []*Role       `gorm:"many2many:xadmin_user_role;association_autoupdate:false;association_autocreate:false" json:"roles"`
	Permissions []*Permission `gorm:"many2many:xadmin_permission_user;association_autoupdate:false;association_autocreate:false;" json:"permissions"`
}

//TableName 用户的表名
func (o User) TableName() string {
	return "xadmin_user"
}

//Role 用户角色
type Role struct {
	DefaultModel
	Name        string       `gorm:"type:varchar(50);" json:"name"`
	Permissions []Permission `gorm:"many2many:xadmin_role_permission;association_autoupdate:false;association_autocreate:false" json:"permissions"`
}

//TableName 用户的表名
func (o Role) TableName() string {
	return "xadmin_role"
}

//UserRole 用户与角色关系
type UserRole struct {
	UserID uint `gorm:"UNIQUE_INDEX:user_role;index"`
	RoleID int  `gorm:"UNIQUE_INDEX:user_role;index"`
}

//TableName 用户的表名
func (o UserRole) TableName() string {
	return "xadmin_user_role"
}

//Permission 权限表
type Permission struct {
	DefaultModel
	Name          string      `gorm:"type:varchar(50);" json:"name"`
	ContentType   ContentType `json:"content_type"`
	ContentTypeID uint        `gorm:"UNIQUE_INDEX:model_code" json:"content_type_id"`
	Code          string      `gorm:"type:varchar(50);UNIQUE_INDEX:model_code" json:"code"` //编码
}

//TableName 权限的表名
func (o Permission) TableName() string {
	return "xadmin_permission"
}

//ContentType 内容类型
type ContentType struct {
	DefaultModel
	AppLabel string `json:"app_label"`
	Model    string `json:"model"`
}

//TableName 内容类型的表名
func (o ContentType) TableName() string {
	return "xadmin_content_type"
}

//RolePermission 组与权限的关系
type RolePermission struct {
	RoleID       int `gorm:"UNIQUE_INDEX:role_permission;index"`
	PermissionID int `gorm:"UNIQUE_INDEX:role_permission;index"`
}

//TableName 用户的表名
func (o RolePermission) TableName() string {
	return "xadmin_role_permission"
}

//PermissionUser 用户与权限关系
type PermissionUser struct {
	PermissionID int  `gorm:"UNIQUE_INDEX:user_permission;index"`
	UserID       uint `gorm:"UNIQUE_INDEX:user_permission;index"`
}

//TableName 用户的表名
func (o PermissionUser) TableName() string {
	return "xadmin_permission_user"
}

//HasPermission 检查是否有权限
func (o *User) HasPermission(perm string) bool {
	return HasPermissionForModel(o, o, perm)
}

//HasPermissionForModel 是否对那个model有权限
func HasPermissionForModel(u *User, model interface{}, perm string) (bl bool) {
	bl = false
	if u.IsSuper {
		bl = true
		return
	}
	ids := make([]uint, 0)
	Db.Model(&PermissionUser{}).Where(PermissionUser{UserID: u.ID}).Pluck("permission_id", &ids)
	rids := make([]uint, 0)
	for _, role := range u.Roles {
		rids = append(rids, role.ID)
	}
	//把角色里的权限查出来
	if len(rids) > 0 {
		ids = append(ids, getPermissionsFromRole(rids)...)
	}

	for _, p := range u.GetPermission() {
		ids = append(ids, p.ID)
	}
	perms := getPermissionsForModel(model, perm)
	for _, id := range ids {
		for _, pe := range perms {
			if id == pe.ID {
				bl = true
				return
			}
		}
	}
	return
}

//getPermissionsFromRole 通过角色取得权限
func getPermissionsFromRole(rids []uint) (ids []uint) {
	Db.Model(&RolePermission{}).Where("role_id in (?)", rids).Pluck("permission_id", &ids)
	return
}

//getPermissions 取得权限
func getPermissionsForModel(model interface{}, perm string) (perms []Permission) {
	ct := GetModelName(model)
	code := GenCodeName(perm, ct.Model)
	Db.Where(&Permission{ContentTypeID: ct.ID, Code: code}).Find(&perms)
	return
}

//Title model 标题
func (o User) Title() string {
	return "用户"
}

//AddRole 添加角色
func AddRole(_, name string) error {
	db := Db.Create(&Role{Name: name})
	return db.Error
}

//AddPermission 添加权限
func AddPermission(model ContentType, code string) error {
	modelname := strings.ToLower(model.Model)
	name := fmt.Sprintf("Can %s %s", code, modelname)
	code = GenCodeName(code, modelname)
	db := Db.FirstOrCreate(&Permission{Code: code, Name: name, ContentType: model}, &Permission{Code: code, ContentTypeID: model.ID})
	return db.Error
}

//GetByUsername 通过用户来查找用户
//guangbiao
func (o *User) GetByUsername(username string) *gorm.DB {
	return Db.First(&o, map[string]interface{}{"username": username})
}

//CheckPassword 检查用户密码
//guangbiao
func (o *User) CheckPassword(password string) bool {
	pass := Cmd5(password, o.Salt)
	return o.Password == pass
}

//UpdateInfo 更新信息
func (o *User) UpdateInfo(info interface{}) *gorm.DB {
	return Db.Model(o).Omit("Roles", "Permissions").Updates(info)
}

//AddUser 添加管理员用户
//guangbiao
func AddUser(username, password string, IsSuper bool) (u User, err error) {
	salt := fmt.Sprintf("%d", time.Now().Unix())
	pass := Cmd5(password, salt)
	u = User{Username: username, Password: pass, Salt: salt, IsSuper: IsSuper}
	db := Db.Create(&u)
	err = db.Error
	return
}

//GetPermission 取得用户的权限
func (o *User) GetPermission() (perms []Permission) {
	pids := make([]int, 0)
	up := PermissionUser{UserID: o.ID}
	Db.Model(&up).Where(up).Pluck("permission_id", &pids)

	roleid := make([]int, 0)
	ur := UserRole{UserID: o.ID}
	Db.Model(&ur).Where(ur).Pluck("role_id", &roleid)

	pids2 := make([]int, 0)
	Db.Model(&RolePermission{}).Where("role_id in (?)", roleid).Pluck("permission_id", &pids2)
	pids = append(pids, pids2...)

	Db.Preload("ContentType").Where("id in (?)", pids).Find(&perms)
	return perms
}

// GetPermissionInfo 取得权限信息
func (o *User) GetPermissionInfo() (perms map[string][]string) {
	ps := o.GetPermission()
	perms = make(map[string][]string)
	for _, p := range ps {
		modelPath := fmt.Sprintf("%s.%s", p.ContentType.AppLabel, p.ContentType.Model)
		if _, ok := perms[modelPath]; !ok {
			perms[modelPath] = make([]string, 0)
		}
		perms[modelPath] = append(perms[modelPath], p.Code)
	}
	return
}

//GetUserByID 通过id获取用户
func (o *User) GetUserByID(id int) *gorm.DB {
	key := id
	if id == 0 {
		key = -1
	}
	return Db.Preload("Roles").Preload("Permissions").First(o, key)
}

//SetPassword SetPassword
func (o *User) SetPassword() {
	if o.Password2 != "" {
		o.Salt = fmt.Sprintf("%d", time.Now().Unix())
		o.Password = Cmd5(o.Password, o.Salt)
	}
}

//AutoMigrate AutoMigrate
func (o *XadminConfig) AutoMigrate() {
	Db.
		AutoMigrate(
			&User{},
			&Role{},
			&UserRole{},
			&Permission{},
			&RolePermission{},
			&PermissionUser{},
			&ContentType{},
		)
}

func (o *XadminConfig) initUser() {
	o.RegisterView(
		Handle{
			Path:   "/login",
			Method: []string{iris.MethodPost},
			Func:   Login,
			Jwt:    false,
		},
		Handle{
			Path:   "/info",
			Method: []string{iris.MethodGet},
			Func:   GetInfo,
			Jwt:    true,
		},
		Handle{
			Path:   "/changepassword",
			Method: []string{iris.MethodPost},
			Func:   ChangePassword,
			Jwt:    true,
		},
	)

	o.Register(&User{}, Config{
		BeforeSave: func(obj interface{}) {
			pointer := reflect.ValueOf(obj)
			m := pointer.MethodByName("SetPassword")
			args := []reflect.Value{}
			m.Call(args)
		},
	})
	o.Register(&Permission{}, Config{})
	o.Register(&Role{}, Config{})
	o.Register(&ContentType{}, Config{})
}
