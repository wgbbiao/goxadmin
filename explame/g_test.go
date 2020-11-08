package main

import (
	"testing"

	"github.com/wgbbiao/goxadmin"
)

func TestMany(t *testing.T) {
	goxadmin.Setdb(DB)
	goxadmin.AddUser("admin", "1234qwer", true)
	// 	// DB.AutoMigrate(&User{}, &CreditCard{})
	// 	var user User
	// 	DB.First(&user, 1)
	// 	user.Name = "123123sdfsdf"

	// 	cc := new(CreditCard)
	// 	cc.ID = 5
	// 	cc.Number = "sdfsdfsdf"

	// 	// user.CreditCards = []CreditCard{
	// 	// 	{Number: "asdfkjas;djfals;kdfjasdfasdfasdfas"},
	// 	// 	*cc,
	// 	// }
	// 	// DB.Create(&User{Name: "ddd", CreditCards: []CreditCard{
	// 	// }})
	// 	// DB.Save(&user)

	// 	DB.Model(&user).Association("CreditCards").Replace([]CreditCard{
	// 		{Number: "asdfkjas;djfals;4"},
	// 		{Number: "asdfkjas;djfals;1"},
	// 		{Number: "asdfkjas;djfals;2"},
	// 		{Number: "asdfkjas;djfals;3"},
	// 		// *cc,
	// 	})

	// 	// sc, _ := schema.Parse(user, &sync.Map{}, DB.NamingStrategy)
	// 	// for _, f := range sc.Relationships.Relations {
	// 	// 	fmt.Println(f.FieldSchema.DBNames) //所有字段
	// 	// 	fmt.Println(f.FieldSchema.Table)   //数据库里的表名
	// 	// 	fmt.Println(f.Name)                //strcut 里的表名
	// 	// 	d := reflect.Indirect(reflect.ValueOf(user))
	// 	// 	ff := d.FieldByName(f.Name)
	// 	// 	px := make([]uint64, 0)
	// 	// 	for i := 0; i < ff.Len(); i++ {
	// 	// 		elm := ff.Index(i)
	// 	// 		// fmt.Println(elm.FieldByName("ID").Uint())
	// 	// 		px = append(px, elm.FieldByName("ID").Uint())
	// 	// 	}
	// 	// 	fmt.Println(px) //关联表所有的ID
	// 	// 	fmt.Println()
	// 	// 	fmt.Println(f.FieldSchema.Relationships.Relations)
	// 	// 	// for _, f := range sc.Fields {
	// 	// 	// 	fmt.Println(f.Name)
	// 	// 	// }
	// 	// }
	// 	// User.id = 1
	// 	// sc, _ := schema.Parse(&CreditCard{}, &sync.Map{}, DB.NamingStrategy)
	// 	// // tableName := "user"
	// 	// // u := sc.LookUpField("User")
	// 	// // fmt.Println(u)
	// 	// if rel, ok := sc.Relationships.Relations["User"]; ok {
	// 	// 	foreignkey := rel.Field.TagSettings["FOREIGNKEY"]
	// 	// 	if foreignkey == "" {
	// 	// 		foreignkey = "UserID"
	// 	// 	}
	// 	// 	fmt.Println(DB.NamingStrategy.ColumnName("User", "UserID"))
	// 	// }

	// 	// for k, f := range sc.Relationships.Relations {
	// 	// 	fmt.Println("----------------")
	// 	// 	fmt.Println(k)
	// 	// 	fmt.Println(f.FieldSchema.DBNames)
	// 	// 	fmt.Println("----------------")
	// 	// 	fmt.Println(f.Field.TagSettings["FOREIGNKEY"])
	// 	// 	fmt.Println(f.References)
	// 	// 	fmt.Println("----------------")
	// 	// }
	// 	fmt.Println("=======")
	// 	// var ccs []CreditCard
	// 	// DB.Joins("User").Where("user.id = ?", 1).Find(&ccs)
	// 	// fmt.Println(DB.NamingStrategy.TableName("user"))
	// 	// fmt.Println(ccs)
	// 	sc, _ := schema.Parse(&CreditCard{}, &sync.Map{}, DB.NamingStrategy)
	// 	fmt.Println(sc.FieldsByDBName)
	// 	fmt.Println(sc.FieldsByName)

}
