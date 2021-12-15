package src

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
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

type myDB sql.DB

var db *myDB
var dbs []database

func (db *myDB) save() {
	d := reflect.ValueOf(db)
	fmt.Println(d.Type())
}

func connect() {
	db = takeDB()

	// доктора
	doctors, err := db.Query("select * from doctors")
	if err != nil {
		panic(err)
	}
	defer func(doctors *sql.Rows) {
		err := doctors.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(doctors)
	var ds []doctor

	for doctors.Next() {
		d := doctor{}
		err := doctors.Scan(&d.sfm, &d.department, &d.specialization, &d.district)
		if err != nil {
			fmt.Println(err.Error() + " doctors")
			continue
		}
		if !d.district.Valid {
			d.district.Int32 = 0
		}
		ds = append(ds, d)
	}
	dbs = append(dbs, ds)

	// пациенты
	patients, err := db.Query("select * from patients")
	if err != nil {
		panic(err)
	}
	var ps []patient

	for patients.Next() {
		p := patient{}
		err := patients.Scan(&p.id, &p.insuranceCompany, &p.sfm, &p.residence, &p.birth, &p.sex, &p.district)
		if err != nil {
			fmt.Println(err.Error() + " patients")
			continue
		}
		ps = append(ps, p)
	}
	dbs = append(dbs, ps)

	// талоны
	talons, err := db.Query("select * from talons")
	if err != nil {
		panic(err)
	}
	defer func(talons *sql.Rows) {
		err := talons.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(talons)
	var ts []talon

	for talons.Next() {
		t := talon{}
		err := talons.Scan(&t.id, &t.reception, &t.doctor, &t.patient)
		if err != nil {
			fmt.Println(err.Error() + " talons")
			continue
		}
		ts = append(ts, t)
	}
	dbs = append(dbs, ts)
}
