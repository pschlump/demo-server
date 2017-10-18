package jwtlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go" // jwt "github.com/dgrijalva/jwt-go"
)

// SignToken Create, sign, and output a token.  This is a great, simple example of
// how to use this library to create and sign a token.
func SignToken(tokData []byte, keyFilePrivate string) (out string, err error) {

	// parse the JSON of the claims
	var claims jwt.MapClaims
	if err = json.Unmarshal(tokData, &claims); err != nil {
		err = fmt.Errorf("Couldn't parse claims JSON: %v", err)
		return
	}

	// get the key
	var key interface{}
	// key, err = loadData(keyFilePrivate)
	key, err = ioutil.ReadFile(keyFilePrivate)
	if err != nil {
		err = fmt.Errorf("Unable to read private key file: %s, err=%v", keyFilePrivate, err)
		return
	}

	// get the signing alg
	alg := jwt.GetSigningMethod("RS256") // Just use RSA256 for now.  May want to make this a parameter.
	if alg == nil {
		err = fmt.Errorf("Couldn't find signing method: %v", "RS256")
		return
	}

	// create a new token
	token := jwt.NewWithClaims(alg, claims)

	if isEs() {
		if k, ok := key.([]byte); !ok {
			err = fmt.Errorf("Couldn't convert key data to key")
			return
		} else {
			key, err = jwt.ParseECPrivateKeyFromPEM(k)
			if err != nil {
				return
			}
		}
	} else if isRs() {
		if k, ok := key.([]byte); !ok {
			err = fmt.Errorf("Couldn't convert key data to key")
			return
		} else {
			key, err = jwt.ParseRSAPrivateKeyFromPEM(k)
			if err != nil {
				return
			}
		}
	}

	if out, err = token.SignedString(key); err != nil {
		err = fmt.Errorf("Error signing token: %v", err)
	}

	return
}

// isEs should return true if we are using an eliptic curve.
func isEs() bool {
	return false
}

// isEs should return true if we are using RSA
func isRs() bool {
	return true
}

/* vim: set noai ts=4 sw=4: */
