package goxadmin

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
)

//User 管理员
type User struct {
	Model
	Username    string        `gorm:"varchar(50);UNIQUE_INDEX" json:"username"`
	Password    string        `gorm:"varchar(50)" json:"password,omitempty"`
	Password2   string        `gorm:"-" json:"password2,omitempty"`
	Salt        string        `gorm:"varchar(64)" json:"-,omitempty"`
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
	Model
	Name        string       `gorm:"varchar(50);" json:"name"`
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
	Model
	Name  string `gorm:"varchar(50);" json:"name"`
	Table string `gorm:"varchar(50);UNIQUE_INDEX:model_code" json:"model"`
	Code  string `gorm:"varchar(10);UNIQUE_INDEX:model_code" json:"code"` //编码
}

//TableName 用户的表名
func (o Permission) TableName() string {
	return "xadmin_permission"
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
	fmt.Println("u.IsSuper == true:", u.IsSuper == true)
	if u.IsSuper == true {
		bl = true
		return
	}
	ids := make([]uint, 0)
	Db.Model(&PermissionUser{}).Where(PermissionUser{UserID: u.ID}).Pluck("permission_id", &ids)
	rids := make([]uint, 0)
	for _, role := range u.Roles {
		rids = append(rids, role.ID)
	}
	if len(rids) > 0 {
		ids = append(ids, getPermissionsFromRole(rids)...)
		fmt.Println("ids:", ids)
		perms := getPermissionsForModel(model, perm)
		fmt.Println("perms:", perms)
		for _, id := range ids {
			for _, pe := range perms {
				if id == pe.ID {
					bl = true
					return
				}
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
	Db.Where(&Permission{Table: GetModelName(model), Code: perm}).Find(&perms)
	return
}

//Title model 标题
func (o User) Title() string {
	return "用户"
}

//AddRole 添加角色
func AddRole(code, name string) error {
	db := Db.Create(&Role{Name: name})
	return db.Error
}

//AddPermission 添加权限
func AddPermission(model, code, name string) error {
	db := Db.FirstOrCreate(&Permission{Code: code, Name: name, Table: model}, &Permission{Code: code, Table: model})
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
func AddUser(username, password string) (u User, err error) {
	salt := fmt.Sprintf("%d", time.Now().Unix())
	pass := Cmd5(password, salt)
	u = User{Username: username, Password: pass, Salt: salt}
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

	// pids3 := make([]int, 0)
	// Db.Table("document_user").Where("user_id = ?", o.ID).Joins("left join permission on document_user.act = permission.code and permission.model=?", "document.Document").Select("permission.id").Pluck("id", &pids3)
	// pids = append(pids, pids3...)

	Db.Where("id in (?)", pids).Find(&perms)
	return perms
}

// GetPermissionInfo 取得权限信息
func (o *User) GetPermissionInfo() (models []string, perms map[string][]string) {
	// ps := o.GetPermission()
	perms = make(map[string][]string)
	models = make([]string, 0)
	// for _, p := range ps {
	// 	if conf.InArrayString(models, p.Model) == false {
	// 		models = append(models, p.Model)
	// 	}
	// 	if _, ok := perms[p.Model]; ok == false {
	// 		perms[p.Model] = make([]string, 0)
	// 	}
	// 	perms[p.Model] = append(perms[p.Model], p.Code)
	// }
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
	o.Salt = fmt.Sprintf("%d", time.Now().Unix())
	o.Password = Cmd5(o.Password, o.Salt)
}

//AutoMigrate AutoMigrate
func AutoMigrate() {
	Db.AutoMigrate(
		&User{},
		&Role{},
		&UserRole{},
		&Permission{},
		&RolePermission{},
		&PermissionUser{},
	)
	Db.Model(&PermissionUser{}).AddForeignKey("user_id", "xadmin_user(id)", "cascade", "cascade")
}

func init() {
	initValidator()
	RegisterView(
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
		})

	Register(&User{}, Config{
		BeforeSave: func(obj interface{}) {
			pointer := reflect.ValueOf(obj)
			m := pointer.MethodByName("SetPassword")
			args := []reflect.Value{}
			m.Call(args)
		},
	})
	Register(&Permission{}, Config{})
	Register(&Role{}, Config{})
}
