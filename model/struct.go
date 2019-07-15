package model

import (
	"github.com/jinzhu/gorm"
)

type (

	Users struct{
		Users []User `json:"users"`
	}
	User struct {
		Name	string `json:"name"`
		Mobile string `json:"mobile"`
	}
	LogModel struct {
		gorm.Model
		Name string `json:"name"`
		Mobile string `json:"mobile"`
	}
	//transformedlogModel struct {
	//	ID uint `json:"id"`
	//	Name	string `json:"name"`
	//	Mobile string `json:"mobile"`
	//}
)
