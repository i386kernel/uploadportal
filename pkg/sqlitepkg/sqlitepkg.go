package sqlitepkg

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
}

type RATable struct {
	AuthorName     string
	Subject        string
	Version        string
	DateTime       string
	LocalFileName  string
	ObjectStoreKey string
}

//func (ra *RATable) CreateDB

func Newsqlite() *RATable {
	return &RATable{}
}

func (rt *RATable) Insert() {

	db, err := sql.Open("sqlite3", "./data/refarch.db")
	if err != nil {
		fmt.Println(err)
	}
	stmt, err := db.Prepare("INSERT INTO ratable(authorname, subject, version, datetime, localfilename, objectstorekey) VALUES (?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(rt.AuthorName, rt.Subject, rt.Version, rt.DateTime, rt.LocalFileName, rt.ObjectStoreKey)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
}

func (rt *RATable) Query() []RATable {

	db, err := sql.Open("sqlite3", "./data/refarch.db")
	if err != nil {
		fmt.Println(err)
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
