package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

type RecordItem struct {
	Id               string
	OriginalUrl      string
	Clicks           int
	CreatedTimestamp int
	Duration         int
}

type Db struct{ *sql.DB }

func Connect(connStr string) *Db {
	fmt.Println("Connecting to the database")
	c, er := sql.Open("postgres", connStr)

	if er != nil {
		log.Fatal(er)
	}

	r, e := c.Exec("SELECT 1 + 1 = 2")

	if e != nil {
		log.Fatal(e)
	}

	rows, _ := r.RowsAffected()
	fmt.Println("Test query succeeded: " + strconv.Itoa(int(rows)))

	return &Db{c}
}

func (db *Db) FindRecordById(id string) (ri *RecordItem, e error) {
	stmt, e := db.Prepare("SELECT * FROM urls WHERE id = $1")

	if e != nil {
		fmt.Println(e.Error())
		return nil, errors.New("could prepare statement")
	}

	defer stmt.Close()

	ri = new(RecordItem)
	e = stmt.QueryRow(id).Scan(&ri.Id, &ri.OriginalUrl, &ri.Clicks, &ri.Duration, &ri.CreatedTimestamp)

	if e != nil {
		fmt.Println(e.Error())
		return nil, errors.New("could prepare statement: ")
	}

	return ri, nil
}

func (db *Db) AddRecord(fi *RecordItem) error {
	stmt, e := db.Prepare("INSERT INTO urls(id, url, clicks, duration, created_timestamp) VALUES($1, $2, $3, $4, $5)")

	if e != nil {
		fmt.Println(e.Error())
		return errors.New("could not prepare statement")
	}

	defer stmt.Close()

	r, er := stmt.Query(fi.Id, fi.OriginalUrl, fi.Clicks, fi.Duration, fi.CreatedTimestamp)

	if er != nil {
		fmt.Println(e.Error())
		return errors.New("could not insert item into db")
	}

	r.Close()
	return nil
}

func (db *Db) IncrementClick(id string) error {
	stmt, e := db.Prepare("UPDATE urls SET clicks = clicks + 1 WHERE id = $1")

	if e != nil {
		fmt.Println(e.Error())
		return errors.New("could not prepare statement")
	}

	defer stmt.Close()

	_, e = stmt.Exec(id)

	if e != nil {
		fmt.Println(e.Error())
		return errors.New("could update item in db")
	}

	return nil
}
