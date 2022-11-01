package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type DBParams struct {
	User     string
	Password string
	Host     string
	Port     string
	Dbname   string
}

type Database struct {
	DB *sql.DB
}

func Init_Database() (*Database, error) {
	log.SetLevel(logrus.DebugLevel)
	config, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Errorf("Failed to fatch config data: %s", err.Error())
		return nil, err
	}
	var dbparams *DBParams
	err = json.Unmarshal(config, &dbparams)
	if err != nil {
		log.Errorf("Faild to unmarshal config file data: %s", err.Error())
		return nil, err
	}
	conf := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbparams.User, dbparams.Password, dbparams.Host, dbparams.Port, dbparams.Dbname)
	db, err := sql.Open("postgres", conf)
	if err != nil {
		log.Errorf("failed to connect database: %s", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Errorf("Failed to ping postgres database: %s", err.Error())
		return nil, err
	}
	d := &Database{DB: db}
	err = d.Users()
	if err != nil {
		log.Errorf("failed to create users table: %s", err.Error())
		return nil, err
	}
	err = d.CreateTable()
	if err != nil {
		log.Errorf("failed to create student table %s", err.Error())
		return nil, err
	}
	return d, nil
}

func (d *Database) CreateTable() error {

	query := `create table if not exists students (
				studentid integer PRIMARY KEY NOT NULL,
				firstname varchar(50) NOT NULL,
				lastname varchar(50) NOT NULL,
				standered varchar(5) NOT NULL,
				percentage decimal NOT NULL,
				isactive boolean default TRUE)`

	_, err := d.DB.Exec(query)
	if err != nil {
		log.Errorf("Failed to create table in psql: %s", err.Error())
		return err
	}
	return nil
}

func (d *Database) Users() error {
	query := `create table if not exists users (
		uid varchar(50) NOT NULL,
		firstname varchar(50) NOT NULL,
		lastname varchar(50) NOT NULL,
		email varchar(50) NOT NULL UNIQUE,
		contact varchar(20) NOT NULL UNIQUE,
		password varchar(100) NOT NULL,
		isactive boolean default true,
		PRIMARY KEY (email,contact))`
	_, err := d.DB.Exec(query)
	if err != nil {
		log.Errorf("faild to create users table %s", err.Error())
		return err
	}
	return nil
}

func (d *Database) SignUp(uid, firstname, lastname, email, contact, hashedpassword string) error {
	query := `insert into users(uid,firstname,lastname,email,contact,password,isactive)Values($1,$2,$3,$4,$5,$6,$7)`
	_, err := d.DB.Exec(query, uid, firstname, lastname, email, contact, hashedpassword, true)
	if err != nil {
		log.Errorf("Failed signup : %s", err.Error())
		return err
	}
	return nil
}

func (d *Database) InsertData(firstname, lastname, studentid, standered, percentage string) error {
	query := "insert into students (studentid,firstname,lastname,standered,percentage) Values($1,$2,$3,$4,$5)"
	val, err := strconv.ParseFloat(percentage, 64)
	if err != nil {
		log.Errorf("Failed to convert into float: %s", err.Error())
		return err
	}
	_, err = d.DB.Exec(query, studentid, firstname, lastname, standered, val)
	if err != nil {
		log.Errorf("Failed to insert data in table: %s", err.Error())
		return err
	}
	return nil
}

func (d *Database) UserLogin(email string) (map[string]string, error) {
	query := "Select password,contact,firstname,lastname from users where email=$1"
	var firstname, lastname, password, contact string
	err := d.DB.QueryRow(query, email).Scan(&password, &contact, &firstname, &lastname)
	if err != nil {
		log.Errorf("failed to facth user details %s", err.Error())
		return nil, err
	}
	mp := make(map[string]string)
	log.Debugf("check out column %s %s %s ", firstname, lastname, contact)
	mp["firstname"] = firstname
	mp["lastname"] = lastname
	mp["password"] = password
	mp["contact"] = contact
	return mp, nil
}
