package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

func fail(msg ...interface{}) {
	fmt.Print("Failed: ")
	fmt.Println(msg...)
	os.Exit(1)
}

func getFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		fail(err)
	}
	return f
}

func main() {

	if len(os.Args) < 2 {
		fail("Missing file path argument")
	}

	inFile := os.Args[1]

	db, err := sql.Open("sqlite3", "./"+inFile+".sqlite")
	if err != nil {
		fail(err)
	}
	defer db.Close()

	db.Exec(`DROP TABLE log`)
	db.Exec(`
		CREATE TABLE log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file VARCHAR(128),
			level VARCHAR(32) NULL,
			process VARCHAR(255) NULL,
			date DATETIME NULL,
			msg TEXT NULL
		);
	`)

	stmt, err := db.Prepare("INSERT INTO log (file, level, process, date, msg) values (?,?,?,?,?)")
	if err != nil {
		fail(err)
	}

	f := getFile(inFile)
	defer f.Close()

	count := 0
	scn := bufio.NewScanner(f)
	for scn.Scan() {

		items := strings.SplitN(strings.Trim(scn.Text(), " "), " ", 5)
		if len(items) != 5 {
			continue
		}

		process := strings.Split(strings.Trim(items[1], "[]"), ":")
		if len(process) != 2 {
			continue
		}

		time := strings.Split(items[3], ",")
		if len(time) != 2 {
			continue
		}

		if _, err := stmt.Exec(inFile, items[0], process[0], items[2]+" "+time[0], items[4]); err != nil {
			fail(err)
		}

		if count%100 == 0 {
			fmt.Print(".")
		}

		if count > 0 && count%10000 == 0 {
			fmt.Print(count)
		}

		count++
	}

	if err := scn.Err(); err != nil {
		fail(err)
	}

	fmt.Println("COMPLETE! ", count, " lines parsed")
}
