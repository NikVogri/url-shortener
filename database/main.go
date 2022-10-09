package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	count := 0

	for {
		conn, e := openConnection(connStr)

		if e != nil {
			log.Println(e.Error())
			count++
		} else {
			fmt.Println("Successfully connected to the database!")

			db := &Db{conn}
			db.createInitialTable()

			return db
		}

		if count > 5 {
			log.Fatal("Could not connect to Postgres database.")
		}

		time.Sleep(5 * time.Second)
	}
}

func openConnection(connStr string) (*sql.DB, error) {
	conn, e := sql.Open("postgres", connStr)

	if e != nil {
		return nil, e
	}

	e = conn.Ping()

	if e != nil {
		return nil, e
	}

	return conn, nil
}

func (db *Db) FindRecordById(id string) (ri *RecordItem, e error) {
	stmt, e := db.Prepare("SELECT * FROM urls WHERE id = $1")

	if e != nil {
		return nil, e
	}

	defer stmt.Close()

	ri = new(RecordItem)
	e = stmt.QueryRow(id).Scan(&ri.Id, &ri.OriginalUrl, &ri.Clicks, &ri.Duration, &ri.CreatedTimestamp)

	if e != nil {
		return nil, e
	}

	return ri, nil
}

func (db *Db) AddRecord(fi *RecordItem) error {
	stmt, e := db.Prepare("INSERT INTO urls(id, url, clicks, duration, created_timestamp) VALUES($1, $2, $3, $4, $5)")

	if e != nil {
		return e
	}

	defer stmt.Close()

	r, er := stmt.Query(fi.Id, fi.OriginalUrl, fi.Clicks, fi.Duration, fi.CreatedTimestamp)

	if er != nil {
		return e
	}

	r.Close()
	return nil
}

func (db *Db) IncrementClick(id string) error {
	stmt, e := db.Prepare("UPDATE urls SET clicks = clicks + 1 WHERE id = $1")

	if e != nil {
		return e
	}

	defer stmt.Close()

	_, e = stmt.Exec(id)

	if e != nil {
		return e
	}

	return nil
}

func (db *Db) createInitialTable() {
	query := "CREATE TABLE IF NOT EXISTS urls ("
	query += "id text PRIMARY KEY, "
	query += "url text NOT NULL, "
	query += "clicks integer NOT NULL, "
	query += "duration bigint NOT NULL, "
	query += "created_timestamp bigint NOT NULL);"

	_, e := db.Exec(query)

	if e != nil {
		log.Panic("Could not create initial table: " + e.Error())
	}
}
