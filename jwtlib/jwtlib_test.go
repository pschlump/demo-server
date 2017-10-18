package jwtlib

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pschlump/demo-server/misc"
)

func Test_JwtTokenLibrary(t *testing.T) {

	if !misc.Exists("../key/sample_key") || !misc.Exists("../key/sample_key.pub") {
		fmt.Printf("Tests require a ../key/sample_key and ../key/sample_key.pub files\n")
		fmt.Printf("These should be self signed RSA256 bit keys\n")
		t.Errorf("Error: Missing require file for tests\n")
	}

	tokData := `{"auth_token":"1234"}`
	jwtSignedToken, err := SignToken([]byte(tokData), "../key/sample_key")

	if err != nil {
		t.Errorf("Error: Error signing token %s\n", err)
	}
	if jwtSignedToken == "" {
		t.Errorf("Error: Error signing token\n")
	}

	authToken, err := VerifyToken([]byte(jwtSignedToken), "../key/sample_key.pub")
	if err != nil {
		t.Errorf("Error: Error signing token %s\n", err)
	}
	if authToken != "1234" {
		t.Errorf("Error: Error signing token\n")
	}

}

func Test_CreateJwtToken(t *testing.T) {

	if !misc.Exists("../key/sample_key") || !misc.Exists("../key/sample_key.pub") {
		fmt.Printf("Tests require a ../key/sample_key and ../key/sample_key.pub files\n")
		fmt.Printf("These should be self signed RSA256 bit keys\n")
		t.Errorf("Error: Missing require file for tests\n")
	}

	tokData := `{"auth_token":"1234"}`

	rv, err := CreateJwtToken(tokData, "../key/sample_key")

	if err != nil {
		t.Errorf("Error: Error signing token %s\n", err)
	}
	if rv == "" {
		t.Errorf("Error: Error signing token\n")
	}
	if !strings.HasPrefix(rv, `{"token":`) {
		t.Errorf("Error: Error signing token\n")
	}

	// fmt.Printf("rv=%s\n", rv)

}

/* vim: set noai ts=4 sw=4: */
