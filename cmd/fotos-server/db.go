package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // MIT licensed.

	"log"
	"sync"
)

var (
	db   *sql.DB
	dbMu sync.Mutex
)

var dblog *log.Logger

var (
	readToken *sql.Stmt

	getTkn  *sql.Stmt
	setAddr *sql.Stmt
	getAddr *sql.Stmt
	setUser *sql.Stmt
)

func openDB(path string) {
	var err error

	dblog = log.New(log.Writer(), "[DB]", 0)
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		dblog.Fatalf("open file %s: %s", path, err)
	}

	// create table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `users` (`name` VARCHAR[255] PRIMARY KEY, `token` VARCHAR[25] KEY, `address` VARCHAR[255]);")
	if err != nil {
		dblog.Fatalf("open file %s: %s", path, err)
	}

	// prepare statements:
	readToken, err = db.Prepare("SELECT name FROM users WHERE token = ?")
	if err != nil {
		dblog.Fatalf("open file %s: %s", path, err)
	}

	getAddr, err = db.Prepare("SELECT address FROM users WHERE name = ?")
	if err != nil {
		dblog.Fatalf("open file %s: %s", path, err)
	}

	getTkn, err = db.Prepare("SELECT token FROM users WHERE name = ?")
	if err != nil {
		dblog.Fatalf("open file %s: %s", path, err)
	}

	setAddr, err = db.Prepare("UPDATE users SET address = ? WHERE token = ?")
	if err != nil {
		dblog.Fatalf("open file %s: %s", path, err)
	}

	setUser, err = db.Prepare("INSERT OR REPLACE INTO users (name, token) VALUES (?, ?)")
	if err != nil {
		dblog.Fatalf("open file %s: %s", path, err)
	}

	return
}

func validateToken(tkn string) bool {
	r := readToken.QueryRow(tkn)

	var name string
	r.Scan(&name)
	if len(name) == 0 {
		return false
	}

	return true
}

func setAddress(tkn, addr string) error {
	_, err := setAddr.Exec(addr, tkn)
	if err != nil {
		dblog.Printf("Tried to set address but an error ocurred: %s\n", err)
		return err
	}

	return nil
}

func getAddress(name string) string {
	r := getAddr.QueryRow(name)

	var addr string
	r.Scan(&addr)
	if len(addr) == 0 {
		dblog.Printf("Query address, but got empty response Name '%s'", name)
	}

	return addr
}

func getToken(name string) string {
	r := getTkn.QueryRow(name)

	var tkn string
	r.Scan(&tkn)
	if len(tkn) == 0 {
		dblog.Printf("Query token, but got empty response Name '%s'", name)
	}

	return tkn
}

func addUser(name, tkn string) error {
	if tkn == "auto" {
		tkn = genTkn()
	}

	_, err := setUser.Exec(name, tkn)
	if err != nil {
		dblog.Printf("Error creating user '%s': %s\n", name, err)
	}

	return err
}
