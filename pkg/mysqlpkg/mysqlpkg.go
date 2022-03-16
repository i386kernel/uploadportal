package mysqlpkg

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Dbcreds struct {
	DBuser string
	DBPass string
	DBName string
	DBsvc  string
}

type RATable struct {
	Dbcreds        Dbcreds
	AuthorName     string
	Subject        string
	Version        string
	DateTime       string
	LocalFileName  string
	ObjectStoreKey string
}

func MewMysqlClient() RATable {
	return RATable{}
}

// InsertRA inserts data to ratable in referencearch database
func (rt *RATable) InsertRA() {
	datasource := fmt.Sprintf("%s:%s%s/%s", rt.Dbcreds.DBuser, rt.Dbcreds.DBPass, rt.Dbcreds.DBsvc, rt.Dbcreds.DBName)
	fmt.Println("Datasource: ", datasource)
	mydbconn, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err.Error())
	}
	stmt, err := mydbconn.Prepare("INSERT INTO ratable(authorName, subject, version, datetime, localfilename, objectstorekey) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(rt.AuthorName, rt.Subject, rt.Version, rt.DateTime, rt.LocalFileName, rt.ObjectStoreKey)
	if err != nil {
		fmt.Println(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(mydbconn)
}

func (rt *RATable) Query() []RATable {

	datasource := fmt.Sprintf("%s:%s%s/%s", rt.Dbcreds.DBuser, rt.Dbcreds.DBPass, rt.Dbcreds.DBsvc, rt.Dbcreds.DBName)
	fmt.Println("Datasource: ", datasource)
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err.Error())
	}
	rows, err := db.Query("SELECT * FROM ratable")
	if err != nil {
		fmt.Println(err)
	}
	rlist := make([]RATable, 0)
	for rows.Next() {
		rtable := RATable{}
		err := rows.Scan(&rtable.AuthorName, &rtable.Subject, &rtable.Version, &rtable.DateTime, &rtable.LocalFileName, &rtable.ObjectStoreKey)
		if err != nil {
			fmt.Println(err)
		}
		rlist = append(rlist, rtable)
	}
	rows.Close()
	return rlist
}
