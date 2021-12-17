package src

import (
	_ "github.com/mattn/go-sqlite3"
)

type database interface{}

type clients struct {
	id        int
	sfm       string
	residence string
}

type librarians struct {
	id  int
	sfm string
}

type films struct {
	name     string
	year     string
	director string
	genre    string
	timeline int
	studio   string
}

type cassettes struct {
	id    int
	price float32
	film  string
	year  int
}

type givings struct {
	id       int
	client   int
	cassette int
	issued   int
}

var dbs []database

func connect() {
	takeDB()
}
