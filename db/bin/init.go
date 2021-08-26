package main

import "github.com/scott-x/myfmt/db"

func main() {
	db.DB.AutoMigrate(&db.GoFile{})
}
