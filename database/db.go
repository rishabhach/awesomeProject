package database

import (
	"awesomeProject/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql","root:rishabh@/logDB?charset=utf8&parseTime=True&loc=Local")
	if err !=nil {
		//panic("failed to connect database")
		panic(err)
	}

	//migrate the schema
	DB.AutoMigrate(&model.LogModel{})
}
