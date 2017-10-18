package jwtlib

import (
	"fmt"

	"github.com/pschlump/demo-server/misc"
	"github.com/pschlump/godebug"
)

// Uses: jwt "github.com/dgrijalva/jwt-go"

// CreateJwtToken will sign a JWT token, with auth_token as a claim inside it.
// func CreateJwtToken(AuthToken string) (rv JwtToken, err error) {
func CreateJwtToken(AuthToken, KeyFilePrivate string) (rv string, err error) {
	all := make(map[string]interface{})

	claims := make(map[string]string)
	claims["auth_token"] = AuthToken
	tokData := godebug.SVar(claims)

	if misc.Db["CreateJwtToken"] {
		fmt.Printf("CreateJWTToken AT:%s: tokData=%s\n", godebug.LF(), tokData)
	}

	signedKey, err := SignToken([]byte(tokData), KeyFilePrivate)
	if err != nil {
		err = fmt.Errorf("Error: Unable to sign the JWT token, %s\n", err)
		return
	}

	all["token"] = signedKey

	rv = godebug.SVar(all)

	if misc.Db["CreateJwtToken"] {
		fmt.Printf("returns rv = %s AT:%s\n", rv, godebug.LF())
	}

	return
}

/* vim: set noai ts=4 sw=4: */
