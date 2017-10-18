package jwtlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go" // jwt "github.com/dgrijalva/jwt-go"
	"github.com/pschlump/demo-server/misc"
	"github.com/pschlump/godebug"
)

// VerifyToken takes a jwt-token and output the auth_token from the claims.  The token's
// signature is verified.  `authToken` is the auth_token in the claims.
func VerifyToken(tokData []byte, keyFile string) (authToken string, err error) {

	// trim possible whitespace from token
	tokData = regexp.MustCompile(`\s*$`).ReplaceAll(tokData, []byte{})
	if misc.Db["db-validate-token"] {
		fmt.Fprintf(os.Stderr, "Token len: %v bytes\n", len(tokData))
	}

	// Parse the token.  Load the key from command line option
	token, err := jwt.Parse(string(tokData), func(t *jwt.Token) (interface{}, error) {
		// data, err := loadData(keyFile)
		data, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		if isEs() {
			return jwt.ParseECPublicKeyFromPEM(data)
		} else if isRs() {
			return jwt.ParseRSAPublicKeyFromPEM(data)
		}
		return data, nil
	})

	// Print some debug data
	if misc.Db["db-validate-token"] && token != nil {
		fmt.Fprintf(os.Stderr, "Header:\n%v\n", token.Header)
		fmt.Fprintf(os.Stderr, "Claims:\n%v\n", token.Claims)
	}

	// Print an error if we can't parse for some reason
	if err != nil {
		return "", fmt.Errorf("Couldn't parse token: %v", err)
	}

	// Is token invalid?
	if !token.Valid {
		return "", fmt.Errorf("Token is invalid")
	}

	if misc.Db["db-token"] {
		fmt.Fprintf(os.Stderr, "Token Claims: %s\n", godebug.SVarI(token.Claims))
	}

	type GetAuthToken struct {
		AuthToken string `json:"auth_token"`
	}
	var gt GetAuthToken
	cl := godebug.SVar(token.Claims)
	if misc.Db["db-jwt-token"] {
		fmt.Fprintf(os.Stderr, "Claims just before -->>%s<<--\n", cl)
	}
	err = json.Unmarshal([]byte(cl), &gt)
	if err != nil {
		if misc.Db["db-jwt-token"] {
			fmt.Fprintf(os.Stderr, "Error: %s -- Unable to unmarsal -->>%s<<--\n", err, cl)
		}
		fmt.Fprintf(os.Stdout, "Error: %s -- Unable to unmarsal -->>%s<<--\n", err, cl)
		return "", err
	}
	if misc.Db["db-jwt-token"] {
		fmt.Fprintf(os.Stderr, "Success: %s -- token [%s] \n", err, gt.AuthToken)
	}
	return gt.AuthToken, nil

}

/* vim: set noai ts=4 sw=4: */
