package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mysqulcrud/internal/filehandler"
	"mysqulcrud/internal/users"
	"mysqulcrud/pkg/miniopkg"
	"mysqulcrud/pkg/mysqlpkg"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"text/template"
)

var mysqltable = mysqlpkg.RATable{}
var minios3object = miniopkg.MinIOObjOptions{}

//var sqlitetable = sqlitepkg.RATable{}

// DefaultHandler handles the first home page request by the client
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	servefile := currentDir + "/static/index.html"
	fmt.Println(servefile)
	http.ServeFile(w, r, servefile)
}

// LoginHandler handles the login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	servefile := currentDir + "/static/login.html"
	fmt.Println(servefile)
	http.ServeFile(w, r, servefile)
}

// ValidateLoginHandler validates the login
func ValidateLoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Printf("Username: %s, Password: %s", username, password)
	if username != "lnanjangud@vmware.com" || password != "password" {
		http.Redirect(w, r, "/incorrectlogin", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// IncorrectLoginHandler handles incorrect login
func IncorrectLoginHandler(w http.ResponseWriter, r *http.Request) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	servefile := currentDir + "/static/incorrectlogin.html"
	fmt.Println(servefile)
	http.ServeFile(w, r, servefile)
}

// SignUpHandler handles signup process
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	users.SignUp(r.FormValue("uname"), r.FormValue("password"))
	servefile := currentDir + "/static/signup.html"
	fmt.Println(servefile)
	http.ServeFile(w, r, servefile)
}

// SignupSuccessHandler handles the signup success
func SignupSuccessHandler(w http.ResponseWriter, r *http.Request) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	servefile := currentDir + "/static/signsucc.html"
	fmt.Println(servefile)
	http.ServeFile(w, r, servefile)
}

// UploadHandler handles the file uploads and reds the data from FormValue
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	flog, err := os.OpenFile("./data/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(flog)
	if runtime.GOOS != "darwin" {
		out, err := exec.Command("cat", "/proc/sys/kernel/hostname").Output()
		if err != nil {
			fmt.Println(err)
		}
		log.Printf("Request Processed in Pod: %s\n", out)
	}
	log.Println("----------------------------------------------------------------------------")
	log.Println(r.FormValue("date"))
	log.Println("----------------------------------------------------------------------------")
	log.Printf("Name: %s || Subject: %s\n", r.FormValue("aname"), r.FormValue("subject"))
	log.Printf("Version: %s || Date/Time: %s\n", r.FormValue("version"), r.FormValue("date"))

	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println(r.FormValue("date"))
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Printf("Name: %s || Subject: %s\n", r.FormValue("aname"), r.FormValue("subject"))
	fmt.Printf("Version: %s || Date/Time: %s\n", r.FormValue("version"), r.FormValue("date"))

	// handle
	fl := filehandler.NewMultiPartHandler(r)
	fl.WriteHTTPMultiPartFormFile()
	if fl.Err != nil {
		http.Error(w, fl.Err.Error(), http.StatusBadRequest)
	}

	// get uploaded filepath and filename
	uploadedfilepath := *fl.Filepath
	uploadedfilename := strings.Split(*fl.Filepath, "/")[2]
	fmt.Printf("Filename: %s\n", uploadedfilename)
	fmt.Printf("Filepath: %s\n", uploadedfilepath)

	// Lookup msql env and populate the appropriate table
	if _, ok := os.LookupEnv("MYSQL_DB_NAME"); ok {
		mysqltable.AuthorName = r.FormValue("aname")
		mysqltable.Subject = r.FormValue("subject")
		mysqltable.Version = r.FormValue("version")
		mysqltable.DateTime = r.FormValue("date")
		mysqltable.ObjectStoreKey = uploadedfilename
		fmt.Println("Inserted Metadata to DB")
		fmt.Println("MYSQL TABLE Object: ", mysqltable)
		mysqltable.InsertRA(mysqlcreds)
	}

	// Lookup S3 bucket name and upload object to minio
	if _, ok := os.LookupEnv("MINIO_BUCKET_NAME"); ok {
		fmt.Println("Object Store Key: ", uploadedfilename)
		minioclient := minIOCreds.NewMinIOClient()
		minios3object.MinClient = minioclient
		minios3object.Location = "default"
		minios3object.Filepath = uploadedfilepath
		minios3object.ObjectName = uploadedfilename
		minios3object.ContentType = "application/mime"
		minios3object.UploadObject()
	}
	http.Redirect(w, r, "/uplsucc", http.StatusSeeOther)
}

// UploadSuccessHandler handles the successful upload
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	servefile := currentDir + "/static/uplsucc.html"
	fmt.Println(servefile)
	http.ServeFile(w, r, servefile)
}

// DownloadHandler handles the download of files
func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	if _, ok := os.LookupEnv("MYSQL_DB_NAME"); ok {
		mysql := mysqlpkg.MewMysqlClient()
		ralist := mysql.Query(mysqlcreds)
		for i, v := range ralist {
			fmt.Println(i, v)
		}
		fp := path.Join("static", "download.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			fmt.Println(err)
		}
		if err := tmpl.Execute(w, ralist); err != nil {
			fmt.Println(err)
			return
		}
	}

}

// SearchHandler searches for existing RA's
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	servefile := currentDir + "/static/search.html"
	fmt.Println(servefile)
}

// LogHandler streams logs upon refresh
func LogHandler(w http.ResponseWriter, r *http.Request) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadFile(currentDir + "/data/app.log")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = fmt.Fprintln(w, string(bs))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// FileServeHandler Serves File dynamically
func FileServeHandler(w http.ResponseWriter, r *http.Request) {
	mv := mux.Vars(r)["filename"]
	fmt.Println("Serving File as per request", mv)
	http.ServeFile(w, r, "./data/"+mv)
}
