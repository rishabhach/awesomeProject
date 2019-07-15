package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"net/http"
	_ "strconv"
	"sync"
)

type (

	Users struct{
		Users []User `json:"users"`
	}
	User struct {
		Name	string `json:"name"`
		Mobile string `json:"mobile"`
	}
	logModel struct {
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

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql","root:rishabh@/logDB?charset=utf8&parseTime=True&loc=Local")
	if err !=nil {
		panic("failed to connect database")
	}

	//migrate the schema
	db.AutoMigrate(&logModel{})
}

func main(){

	fmt.Println("rishabh acharya")
	//http.HandleFunc("/upload", uploadFile)
	//http.ListenAndServe(":8080", nil)
	router := gin.Default()

	v1 :=router.Group("logs")
	{
		v1.POST("/",createLog)
		v1.POST("/upload",uploadFile)
		v1.GET("/admin/:searchId",adminLogs)
		v1.GET("/user/:searchId",userLogs)
	}
	router.Run()
}

func adminLogs(c *gin.Context){
	searchId := c.Param("searchId")
	var users []logModel
	var tuser []User
	db.Where("name=?",searchId).Find(&users)
	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	var waitgroup sync.WaitGroup
	waitgroup.Add(len(users))
	//transforms the todos for building a good response
	for i:=0;i<len(users);i++ {
		var user logModel
		user =users[i]
		fmt.Println(user)
		go func (user *logModel,waitgroup *sync.WaitGroup,tuser *[]User) {
			*tuser = append(*tuser, User{Name: user.Name, Mobile: user.Mobile})
			waitgroup.Done()
			fmt.Println(tuser)
		}(&user,&waitgroup,&tuser)
		//_tuser = append(_tuser, User{Name: user.Name, Mobile: user.Mobile})
		waitgroup.Wait()
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": tuser})
}

func userLogs(c *gin.Context){
	searchId := c.Param("searchId")
	var users []logModel
	var _tuser []User
	db.Where("name=?",searchId).Find(&users)
	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}


	//transforms the todos for building a good response


	var waitgroup sync.WaitGroup
	waitgroup.Add(len(users))


	//transforms the todos for building a good response


	for i:=0;i<len(users);i++ {
		var user logModel
		user =users[i]
		go func (user *logModel,waitgroup *sync.WaitGroup,_tuser *[]User) {
			Maskedmobile := "XXXXXXXXXX"
			*_tuser = append(*_tuser, User{Name: user.Name, Mobile: Maskedmobile})
			waitgroup.Done()
		}(&user,&waitgroup,&_tuser)
		//_tuser = append(_tuser, User{Name: user.Name, Mobile: user.Mobile})
		waitgroup.Wait()
	}
	//for i:=0;i<len(users);i++ {
	//
	//	user :=users[i]
	//	Maskedmobile := "XXXXXXXXXX"
	//	_tuser = append(_tuser, User{Name: user.Name, Mobile: Maskedmobile})
	//}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _tuser})
}

func uploadFile(c *gin.Context){
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	c.Request.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("/Users/rishabhacharya/go/src/temp-images", "upload-*.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fmt.Println(file," file to check")
	fileBytes, err := ioutil.ReadAll(file)


	if err != nil {
		fmt.Println(err)
	}
	var users Users

	json.Unmarshal(fileBytes, &users)
	fmt.Println(users," users to check")

	var user User
	fmt.Println(users.Users,"length to check")

	var waitgroup sync.WaitGroup
	waitgroup.Add(len(users.Users))


	for i:=0;i< len(users.Users);i++ {

		user = users.Users[i]
		addlog(&user,&waitgroup)


		//todo := logModel{Name: user.Name, Mobile: user.Mobile}
		//fmt.Println(user.Name,user.Mobile)
		//db.Save(&todo)
	}
	waitgroup.Wait()
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)



	// return that we have successfully uploaded our file!
	//fmt.Fprintf(w, "Successfully Uploaded File\n")
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Successfully Uploaded File\n"})
	//fmt.Fprintf(w, "Uploading File")
}

func addlog(user *User,waitgroup *sync.WaitGroup){
	todo := logModel{Name: user.Name, Mobile: user.Mobile}
	fmt.Println(user.Name,user.Mobile)
	db.Save(&todo)
	waitgroup.Done()
}

// createTodo add a new todo
func createLog(c *gin.Context) {
	//completed, _ := strconv.Atoi(c.PostForm("completed"))
	logEntry := logModel{Name: c.PostForm("name"), Mobile: c.PostForm("mobile")}
	db.Save(&logEntry)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": logEntry.ID})
}