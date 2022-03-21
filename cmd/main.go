package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"mysqulcrud/pkg/miniopkg"
	"mysqulcrud/pkg/mysqlpkg"
	"net/http"
	"os"
	"time"
)

var mysqlcreds mysqlpkg.Dbcreds
var minIOCreds miniopkg.MinIOCreds

// Env Lookups, instantiate required objects
func init() {

	// Creating necessary directories
	err := os.MkdirAll("./data", os.ModePerm)
	if err != nil {
		panic(err)
	}

	if _, ok := os.LookupEnv("MYSQL_DB_NAME"); ok {
		fmt.Println("MYSQL config found, Waiting for Mysql Service to be ready")
		mysqlcreds.DBuser = os.Getenv("MYSQL_USER")
		mysqlcreds.DBPass = os.Getenv("MYSQL_ROOT_PASSWORD")
		mysqlcreds.DBName = os.Getenv("MYSQL_DB_NAME")
		mysqlcreds.DBsvc = os.Getenv("MYSQL_SVC")
		fmt.Println("mysqlcreds: ", mysqlcreds)
		time.Sleep(1 * time.Minute)
		fmt.Println("Setting up Database")
		mysqlcreds.DBsvc = fmt.Sprintf("@tcp(%s:3306)", mysqlcreds.DBsvc)
		datasource := fmt.Sprintf("%s:%s%s/%s", mysqlcreds.DBuser, mysqlcreds.DBPass, mysqlcreds.DBsvc, "")
		fmt.Println("Datasource: ", datasource)

		mydbconn, err := sql.Open("mysql", datasource)
		if err != nil {
			panic(err.Error())
		}
		// setting up database
		_, err = mydbconn.Exec("CREATE DATABASE IF NOT EXISTS " + mysqlcreds.DBName)
		if err != nil {
			fmt.Println(err)
		}
		// creating schema
		_, err = mydbconn.Exec("USE " + mysqlcreds.DBName)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Running Init, Instantiating Objects")
		fmt.Println("Reading ENV's, Grabbing Values")
		_, err = mydbconn.Exec("CREATE TABLE IF NOT EXISTS ratable(authorname varchar(255), subject varchar(255), version varchar(255), datetime varchar(255), localfilename varchar(255), objectstorekey varchar(1000))")
		if err != nil {
			fmt.Println(err)
		}
	}
	if _, ok := os.LookupEnv("MINIO_BUCKET_NAME"); ok {
		fmt.Println("MINIO config found, accessing MINIO")
		minIOCreds.Endpoint = os.Getenv("MINIO_ENDPOINT")
		minIOCreds.AccessKeyID = os.Getenv("MINIO_ROOT_USER")
		minIOCreds.SecretAccessKey = os.Getenv("MINIO_ROOT_PASSWORD")
		minios3object.Bucketname = os.Getenv("MINIO_BUCKET_NAME")
		fmt.Println("minIOCreds:", minIOCreds)
	} else {
		fmt.Println("MINIO not detected, skipping")
	}

	if _, ok := os.LookupEnv("REDIS_HOST"); ok {
		fmt.Println("REDIS Config Found, yet to be implemented")
	} else {
		fmt.Println("REDIS config not found, maintaining sessions locally, single instance, " +
			"you won't be able to scale")
	}
}

// RouteHandler handles all the routes that routes to appropriate end-points
func RouteHandler() {
	port := ":8080"
	r := mux.NewRouter()
	r.HandleFunc("/", LoginHandler)
	r.HandleFunc("/vallogin", ValidateLoginHandler)
	r.HandleFunc("/incorrectlogin", IncorrectLoginHandler)
	r.HandleFunc("/home", DefaultHandler)
	r.HandleFunc("/upload", UploadHandler)
	r.HandleFunc("/signup", SignUpHandler)
	r.HandleFunc("/signsucc", SignupSuccessHandler)
	r.HandleFunc("/uplsucc", UploadSuccessHandler)
	r.HandleFunc("/download", DownloadHandler)
	r.HandleFunc("/search", SearchHandler)
	r.HandleFunc("/log", LogHandler)
	r.HandleFunc("/downloadfile/{filename}", FileServeHandler)
	fmt.Println("Serving Port = ", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

func main() {
	RouteHandler()
}
