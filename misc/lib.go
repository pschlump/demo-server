package misc

import (
	"errors"
	"fmt"
	"os"

	"github.com/pschlump/json" //	"encoding/json" modified to ignore invalid types.
)

// SVarI converts data to JSON.
func SVar(v interface{}) string {
	s, err := json.Marshal(v)
	// s, err := json.MarshalIndent ( v, "", "\t" )
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}

// SVarI converts data to JSON with indentation.
func SVarI(v interface{}) string {
	// s, err := json.Marshal ( v )
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}

// Exists returns true if file exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Used in Fopen, if mode is not "r", "w", "a" this is returned.
var invalidMode = errors.New("Invalid Mode")

// Fopen opens a file and returns the *os.File.
func Fopen(fn string, mode string) (file *os.File, err error) {
	file = nil
	if mode == "r" {
		file, err = os.Open(fn) // For read access.
	} else if mode == "w" {
		file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	} else if mode == "a" {
		file, err = os.OpenFile(fn, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		}
	} else {
		err = invalidMode
	}
	return
}

/* vim: set noai ts=4 sw=4: */
