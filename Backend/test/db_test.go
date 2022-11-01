package test

import (
	"log"
	"testing"

	db "github.com/UnivertsityStudent/backend/database"
)

func TestDBConnection(t *testing.T) {
	database, err := db.Init_Database()
	if err != nil {
		t.Fatalf("failed init db %s", err.Error())
		return
	}
	database.CreateTable()
	if err != nil {
		log.Fatalf("Failed to create table in psql %s", err.Error())
		return
	}
	t.Logf("check error ")
}

func TestInsertData(t *testing.T) {
	database, err := db.Init_Database()
	if err != nil {
		t.Fatalf("failed init db %s", err.Error())
		return
	}

	err = database.InsertData("Rajkumar", " Raigar", "12334", "5th", "70.00")
	if err != nil {
		t.Fatalf("Failed to insert data in student table: %s", err.Error())
		return
	}
	t.Logf("Successs the test")
}
