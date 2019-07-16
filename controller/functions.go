package controller

import (
	"awesomeProject/database"
	"awesomeProject/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	_ "strconv"
	"sync"
)

func AdminLogs(c *gin.Context){
	searchId := c.Param("searchId")
	var users []model.LogModel
	var tuser []model.User


	database.DB.Where("name=?",searchId).Find(&users)
	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No log found!"})
		return
	}

	var waitgroup sync.WaitGroup
	waitgroup.Add(len(users))

	
	for i:=0;i<len(users);i++ {
		var user model.LogModel
		user =users[i]
		fmt.Println(user)
		go func (user *model.LogModel,waitgroup *sync.WaitGroup,tuser *[]model.User) {
			*tuser = append(*tuser, model.User{Name: user.Name, Mobile: user.Mobile})

			fmt.Println(tuser, "i am in loop ",i)
			waitgroup.Done()
		}(&user,&waitgroup,&tuser)
		//_tuser = append(_tuser, User{Name: user.Name, Mobile: user.Mobile})

	}
	waitgroup.Wait()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": tuser})
}

func UserLogs(c *gin.Context){
	searchId := c.Param("searchId")
	var users []model.LogModel
	var _tuser []model.User
	database.DB.Where("name=?",searchId).Find(&users)
	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No log found!"})
		return
	}


	


	var waitgroup sync.WaitGroup
	waitgroup.Add(len(users))


	


	for i:=0;i<len(users);i++ {
		var user model.LogModel
		user =users[i]
		go func (user *model.LogModel,waitgroup *sync.WaitGroup,_tuser *[]model.User) {
			Maskedmobile := "XXXXXXXXXX"
			*_tuser = append(*_tuser, model.User{Name: user.Name, Mobile: Maskedmobile})
			waitgroup.Done()
		}(&user,&waitgroup,&_tuser)
		//_tuser = append(_tuser, User{Name: user.Name, Mobile: user.Mobile})

	}
	waitgroup.Wait()
	
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _tuser})
}

func UploadFileStruct(c *gin.Context){
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
	var users model.Users

	json.Unmarshal(fileBytes, &users)
	fmt.Println(users," users to check")

	var user model.User
	fmt.Println(users.Users,"length to check")

	var waitgroup sync.WaitGroup
	waitgroup.Add(len(users.Users))


	for i:=0;i< len(users.Users);i++ {

		user = users.Users[i]
		go addlog(&user,&waitgroup)


		
	}
	waitgroup.Wait()
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)



	// return that we have successfully uploaded our file!
	//fmt.Fprintf(w, "Successfully Uploaded File\n")
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Successfully Uploaded File\n"})
	//fmt.Fprintf(w, "Uploading File")
}

func addlog(user *model.User,waitgroup *sync.WaitGroup){
	todo := model.LogModel{Name: user.Name, Mobile: user.Mobile}
	defer waitgroup.Done()
	fmt.Println(user.Name,user.Mobile)
	database.DB.Save(&todo)

}

func UploadFileStructMultiFile(c *gin.Context){
	/*
	 read all files
	 for all files {
	    wg waitgroup
	    wg.Add(1)
	    go func(file) {
	      process file ()
	      wg.Done()
		}(file)
	}
	wg.wait()

	*/
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	c.Request.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file

	fhs := c.Request.MultipartForm.File["myFile"]

	var wg sync.WaitGroup
	wg.Add(len(fhs))


	for i:=0;i<len(fhs);i++ {

		file,err := fhs[i].Open()

		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		fmt.Printf("%T file type\n",file)
		//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		//fmt.Printf("File Size: %+v\n", handler.Size)
		//fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern

		go func(file *multipart.File,wg *sync.WaitGroup){

			fmt.Println("i in first for loop")
			tempFile, err := ioutil.TempFile("/Users/rishabhacharya/go/src/temp-images", "upload-*.json")
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()

			// read all of the contents of our uploaded file into a
			// byte array
			fmt.Println(file," file to check")
			fileBytes, err := ioutil.ReadAll(*file)


			if err != nil {
				fmt.Println(err)
			}
			var users model.Users

			json.Unmarshal(fileBytes, &users)
			fmt.Println(users," users to check")


			fmt.Println(users.Users,"length to check")

			var waitgroup sync.WaitGroup
			waitgroup.Add(len(users.Users))


			for i:=0;i< len(users.Users);i++ {
				var user model.User
				user = users.Users[i]
				go addlog(&user,&waitgroup)


				//todo := logModel{Name: user.Name, Mobile: user.Mobile}
				//fmt.Println(user.Name,user.Mobile)
				//db.Save(&todo)
				fmt.Println("yes i am here in loop d concurrent")
			}
			waitgroup.Wait()
			// write this byte array to our temporary file
			tempFile.Write(fileBytes)
			fmt.Println("i out of second for loop")
			wg.Done()
		}(&file,&wg)




	}
	wg.Wait()

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Successfully Uploaded File\n"})

	//file, handler, err := c.Request.FormFile("myFile")
	//if err != nil {
	//	fmt.Println("Error Retrieving the File")
	//	fmt.Println(err)
	//	return
	//}
	//defer file.Close()
	//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	//fmt.Printf("File Size: %+v\n", handler.Size)
	//fmt.Printf("MIME Header: %+v\n", handler.Header)
	//
	//// Create a temporary file within our temp-images directory that follows
	//// a particular naming pattern
	//tempFile, err := ioutil.TempFile("/Users/rishabhacharya/go/src/temp-images", "upload-*.json")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer tempFile.Close()
	//
	//// read all of the contents of our uploaded file into a
	//// byte array
	//fmt.Println(file," file to check")
	//fileBytes, err := ioutil.ReadAll(file)
	//
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//var users model.Users
	//
	//json.Unmarshal(fileBytes, &users)
	//fmt.Println(users," users to check")
	//
	//var user model.User
	//fmt.Println(users.Users,"length to check")
	//
	//var waitgroup sync.WaitGroup
	//waitgroup.Add(len(users.Users))
	//
	//
	//for i:=0;i< len(users.Users);i++ {
	//
	//	user = users.Users[i]
	//	addlog(&user,&waitgroup)
	//
	//
	//	//todo := logModel{Name: user.Name, Mobile: user.Mobile}
	//	//fmt.Println(user.Name,user.Mobile)
	//	//db.Save(&todo)
	//}
	//waitgroup.Wait()
	//// write this byte array to our temporary file
	//tempFile.Write(fileBytes)
	//
	//

	// return that we have successfully uploaded our file!
	//fmt.Fprintf(w, "Successfully Uploaded File\n")
	//c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Successfully Uploaded File\n"})
	//fmt.Fprintf(w, "Uploading File")
}

/*
func UploadFileUnstruct(c *gin.Context){
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

//	var users model.Users
	var result map[string]interface{}
	json.Unmarshal(fileBytes, &result)
	fmt.Println(result," users to check")

	//var user model.User
	//fmt.Println(users.Users,"length to check")
	//
	//var waitgroup sync.WaitGroup
	//waitgroup.Add(len(users.Users))
	//
	//
	//for i:=0;i< len(users.Users);i++ {
	//
	//	user = users.Users[i]
	//	addlog(&user,&waitgroup)
	//
	//
	//	//todo := logModel{Name: user.Name, Mobile: user.Mobile}
	//	//fmt.Println(user.Name,user.Mobile)
	//	//db.Save(&todo)
	//}
	//waitgroup.Wait()
	// write this byte array to our temporary file

	fmt.Println("yes it is working")
	tempFile.Write(fileBytes)



	// return that we have successfully uploaded our file!
	//fmt.Fprintf(w, "Successfully Uploaded File\n")
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Successfully Uploaded File\n"})
	//fmt.Fprintf(w, "Uploading File")

}*/

// createTodo add a new todo
func CreateLog(c *gin.Context) {
	//completed, _ := strconv.Atoi(c.PostForm("completed"))
	logEntry := model.LogModel{Name: c.PostForm("name"), Mobile: c.PostForm("mobile")}
	database.DB.Save(&logEntry)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": logEntry.ID})
}
