package main

import (
	"awesomeProject/controller"

	_ "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "io/ioutil"
	_ "net/http"
	_ "strconv"
	_ "sync"
)


func main(){

	fmt.Println("rishabh acharya")
	//http.HandleFunc("/upload", uploadFile)
	//http.ListenAndServe(":8080", nil)
	router := gin.Default()

	v1 :=router.Group("logs")
	{
		v1.POST("/", controller.CreateLog)
		v1.POST("/upload", controller.UploadFile)
		v1.GET("/admin/:searchId", controller.AdminLogs)
		v1.GET("/user/:searchId", controller.UserLogs)
	}
	router.Run()
}

