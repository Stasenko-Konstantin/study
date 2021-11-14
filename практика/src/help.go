package src

import (
	"errors"
	"log"
	"os"
)

type myLogger struct {
	file   *os.File
	stdout *os.File
}

func (ml myLogger) Write(b []byte) (n int, err error) {
	var myerr string
	n, err = ml.file.Write(b)
	if err != nil {
		myerr += err.Error() + "\n"
	}
	_, err = ml.stdout.Write(b)
	if err != nil {
		myerr += err.Error()
	}
	if myerr != "" {
		err = errors.New(myerr)
	}
	return n, err
}

func newLogger() *log.Logger {
	if _, err := os.Stat("log.txt"); os.IsNotExist(err) {
		f, err := os.Create("log.txt")
		if err != nil {
			panic("Не удалось создать логер!")
		}
		return log.New(myLogger{f, os.Stdout}, "INFO\t", log.Ltime)
	}
	f, err := os.Open("log.txt")
	if err != nil {
		panic("Не удалось открыть логер!")
	}
	return log.New(myLogger{f, os.Stdout}, "INFO\t", log.Ltime)
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func stringConcat(vals []string) string {
	var r string
	for _, e := range vals {
		r += e + " "
	}
	return r
}

func contain(mas []string, e string) bool {
	for _, s := range mas {
		if s == e {
			return true
		}
	}
	return false
}

func makeDocs(doctors [][]string) []string {
	var r []string
	i := 0
	for _, d := range doctors {
		if i == 0 {
			i += 1
		} else {
			if d[3] != "0" {
				r = append(r, d[0])
			}
		}
	}
	return r
}

func findDoctor(doc string, doctors [][]string) string {
	var r string
	for _, d := range doctors {
		if doc == d[0] {
			r = d[3]
		}
	}
	return r
}

//func (rs *Rows) ScanString(dest *[]*NullString) error {
//	rs.closemu.RLock()
//
//	if rs.lasterr != nil && rs.lasterr != io.EOF {
//		rs.closemu.RUnlock()
//		return rs.lasterr
//	}
//	if rs.closed {
//		err := rs.lasterrOrErrLocked(errRowsClosed)
//		rs.closemu.RUnlock()
//		return err
//	}
//	rs.closemu.RUnlock()
//
//	if rs.lastcols == nil {
//		return errors.New("sql: Scan called without calling Next")
//	}
//	if len(*dest) != len(rs.lastcols) {
//		for i := 0; i < len(rs.lastcols); i++ {
//			var s NullString
//			*dest = append(*dest, &s)
//		}
//	}
//	for i, sv := range rs.lastcols {
//		e := *dest
//		err := convertAssignRows(e[i], sv, rs)
//		if err != nil {
//			return fmt.Errorf(`sql: Scan error on column index %d, name %q: %w`, i, rs.rowsi.Columns()[i], err)
//		}
//	}
//	return nil
//}
