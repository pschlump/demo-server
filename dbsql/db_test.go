package dbsql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/lib/pq"
	JsonX "github.com/pschlump/JSONx"
	"github.com/pschlump/godebug"
)

var hasRun = false

var db *sql.DB

type CfgType struct {
	PGConn string `gfJsonX:"PGConn"`
}

var Cfg CfgType

func SetupForTest(fn string) {
	if !hasRun {

		In, err := ioutil.ReadFile(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test-code: Error returned from ReadFile: %s\n", err)
			os.Exit(1)
		}

		_, err = JsonX.Unmarshal(fn, In, &Cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test-code: Error returned from JsonX.Unmarshal: %s\n", err)
			os.Exit(1)
		}

		auth := Cfg.PGConn

		db_x := ConnectToAnyDb("postgres", auth, "")
		db = db_x.Db
	}
	hasRun = true
}

type ErrType int

const (
	ErrNone    ErrType = 1
	ErrHaveErr ErrType = 2
)

func Test_SelData2(t *testing.T) {

	SetupForTest("../cfg.jsonx")

	tests := []struct {
		stmt    string        //
		data    []interface{} //
		errFlag ErrType
		expRows int
		expData string
	}{
		{ // Test really simple select, calls SelQ and RowsToInterface when SelData2 is called.
			stmt:    "select 12 as \"col\"",
			data:    []interface{}{},
			errFlag: ErrNone,
			expRows: 1,
			expData: `[{"col":12}]`,
		},
		{ // Test passing bind variables
			stmt:    "select 12 as \"col\" where $1 < $2 ",
			data:    []interface{}{"a", "b"},
			errFlag: ErrNone,
			expRows: 1,
			expData: `[{"col":12}]`,
		},
		{ // Test passing bind variables
			stmt:    "select 12 as \"col\" where $1 < $2 ",
			data:    []interface{}{1, 2},
			errFlag: ErrNone,
			expRows: 1,
			expData: `[{"col":12}]`,
		},
		{ // Test passing bind variables, empty result
			stmt:    "select 12 as \"col\" where $1 < $2 ",
			data:    []interface{}{2, 1},
			errFlag: ErrNone,
			expRows: 0,
			expData: `null`, // note null is returned instead of [] or {}
		},
	}

	for ii, test := range tests {

		// func SelData2(db *sql.DB, q string, data ...interface{}) ([]map[string]interface{}, error) {
		mdata, err := SelData2(db, test.stmt, test.data...)
		if err != nil && ii == 0 {
			fmt.Printf(`  If the error you are seeing is about connection to the database then your
  connect string is not correct.  If you are using {{ __env__ DB_PASS }} in the ../cfg.jsonx
  file you many not have exported the password.  Try
    $ export DB_PASS=<<<your password>>>
	$ go test
`)
		}
		if test.errFlag == ErrNone && err != nil {
			t.Errorf(fmt.Sprintf("Error %2d, stmt: %s, bind vars=%v\n", ii, test.stmt, test.data))
		}
		if test.errFlag == ErrHaveErr && err == nil {
			t.Errorf(fmt.Sprintf("Error %2d, stmt: %s, bind vars=%v\n", ii, test.stmt, test.data))
		}
		if test.errFlag == ErrNone && err == nil {
			if len(mdata) != test.expRows {
				t.Errorf(fmt.Sprintf("Error %2d, stmt: %s, bind vars=%v, number of rows got %d, expected %d\n", ii, test.stmt, test.data, len(mdata), test.expRows))
			} else {
				sd := godebug.SVar(mdata)
				if sd != test.expData {
					t.Errorf(fmt.Sprintf("Error %2d, stmt: %s, bind vars=%v, got ->%s<- expected ->%s<-\n", ii, test.stmt, test.data, sd, test.expData))
				}
			}
		}

	}

}

func Test_Run1(t *testing.T) {

	SetupForTest("../cfg.jsonx")

	tests := []struct {
		stmt    string        //
		data    []interface{} //
		errFlag ErrType
	}{
		{ // Test really simple select.
			stmt:    "select 12 as \"col\"",
			data:    []interface{}{},
			errFlag: ErrNone,
		},
		{ // Test passing bind variables
			stmt:    "select 12 as \"col\" where $1 < $2 ",
			data:    []interface{}{"a", "b"},
			errFlag: ErrNone,
		},
		{ // Test passing bind variables
			stmt:    "select 12 as \"col\" where $1 < $2 ",
			data:    []interface{}{1, 2},
			errFlag: ErrNone,
		},
		{ // Test passing bind variables, empty result
			stmt:    "select 12 as \"col\" where $1 < $2 ",
			data:    []interface{}{2, 1},
			errFlag: ErrNone,
		},
	}

	for ii, test := range tests {

		//func Run1(db *sql.DB, q string, arg ...interface{}) error {
		err := Run1(db, test.stmt, test.data...)
		if test.errFlag == ErrNone && err != nil {
			t.Errorf(fmt.Sprintf("Error %2d, stmt: %s, bind vars=%v\n", ii, test.stmt, test.data))
		}
		if test.errFlag == ErrHaveErr && err == nil {
			t.Errorf(fmt.Sprintf("Error %2d, stmt: %s, bind vars=%v\n", ii, test.stmt, test.data))
		}

	}

}

/* vim: set noai ts=4 sw=4: */
