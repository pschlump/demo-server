package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	JsonX "github.com/pschlump/JSONx"
	"github.com/pschlump/demo-server/dbsql"
	"github.com/pschlump/godebug"
)

type TokenType struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}

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

		db_x := dbsql.ConnectToAnyDb("postgres", auth, "")
		db = db_x.Db
	}
	hasRun = true
}

func main() {

	errCount := 0
	var tt TokenType

	// -----------------------------------------------------------------------------------------------------------------------
	// Connect to d.b. - Clean out xyz@test user

	SetupForTest("../cfg.jsonx")

	stmt := `delete from "az_auth_token" where "user_id" in ( select "id" from "az_user" where "email" = 'xyz@test-example....' )`
	err := dbsql.Run1(db, stmt)
	if err != nil {
		fmt.Printf("Error during init (may not be fatal), err=%s\n", err)
		fmt.Printf(`  If the error you are seeing is about connection to the database then your
  connect string is not correct.  If you are using {{ __env__ DB_PASS }} in the ../cfg.jsonx
  file you many not have exported the password.  Try
    $ export DB_PASS=<<<your password>>>
	$ go run testApi.go
`)
	}
	stmt = `delete from "az_user" where "email" = 'xyz@test-example....'`
	err = dbsql.Run1(db, stmt)
	if err != nil {
		fmt.Printf("FAIL\nFatal error unable to setup database for test, err=%s\n", err)
		os.Exit(1)
	}

	// -----------------------------------------------------------------------------------------------------------------------
	// Do a Signup
	{
		client := &http.Client{}
		postData := map[string]string{"email": "xyz@test-example....", "password": "password1", "firstName": "aa first", "lastName": "aa last"}
		jsonValue, _ := json.Marshal(postData)
		req, err := http.NewRequest("POST", "http://localhost:18000/signup", bytes.NewReader(jsonValue))
		req.Header.Add("User-Agent", "Test API Client")
		req.Header.Add("Content-Type", "application/json") // set content type/encoding to JSON
		resp, err := client.Do(req)
		if err != nil {
			errCount++
			fmt.Printf("Error: HTTP failure, %s\n", err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errCount++
				fmt.Printf("Error: failed to read body of responce, %s\n", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					errCount++
					fmt.Printf("Status not 200 on /signup: status=%d, AT:%s\n", resp.StatusCode, godebug.LF())
				} else {
					if db1 || errCount > 0 {
						fmt.Printf("body ->%s<-\n", body)
					}
					err = json.Unmarshal(body, &tt)
					if err != nil {
						errCount++
						fmt.Printf("Error: failed to parse return JSON data, err=%s AT:%s data= ->%s<-\n", err, godebug.LF(), body)
					}
					if len(tt.Token) < 50 {
						errCount++
						fmt.Printf("Error: token is too short, AT:%s data= ->%s<-\n", godebug.LF(), body)
					}
				}
			}
		}
	}

	// -----------------------------------------------------------------------------------------------------------------------
	// do Signup again with same account, sholuld fail.
	{
		client := &http.Client{}
		postData := map[string]string{"email": "xyz@test-example....", "password": "password1", "firstName": "aa first", "lastName": "aa last"}
		jsonValue, _ := json.Marshal(postData)
		req, err := http.NewRequest("POST", "http://localhost:18000/signup", bytes.NewReader(jsonValue))
		req.Header.Add("User-Agent", "Test API Client")
		req.Header.Add("Content-Type", "application/json") // set content type/encoding to JSON
		resp, err := client.Do(req)
		if err != nil {
			errCount++
			fmt.Printf("Error: HTTP failure, %s\n", err)
		} else {
			if resp.StatusCode == 200 {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					errCount++
					fmt.Printf("Error: failed to read body of responce, %s\n", err)
				} else {
					defer resp.Body.Close()
					if db1 || errCount > 0 {
						fmt.Printf("body ->%s<-\n", body)
					}
					err = json.Unmarshal(body, &tt)
					if err != nil {
						errCount++
						fmt.Printf("Error: failed to parse return JSON data, err=%s AT:%s data= ->%s<-\n", err, godebug.LF(), body)
					}
					if tt.Status == "succes" || tt.Status == "" {
						errCount++
						fmt.Printf("Error: failed to get error back, AT:%s data= ->%s<-\n", godebug.LF(), body)
					}
				}
			} else if resp.StatusCode != 400 {
				errCount++
				fmt.Printf("Status not 400 on /signup: status=%d, AT:%s\n", resp.StatusCode, godebug.LF())
			}
		}
	}

	// -----------------------------------------------------------------------------------------------------------------------
	// do Login - should get token.
	{
		client := &http.Client{}
		postData := map[string]string{"email": "xyz@test-example....", "password": "password1"}
		jsonValue, _ := json.Marshal(postData)
		req, err := http.NewRequest("POST", "http://localhost:18000/login", bytes.NewReader(jsonValue))
		req.Header.Add("User-Agent", "Test API Client")
		req.Header.Add("Content-Type", "application/json") // set content type/encoding to JSON
		resp, err := client.Do(req)
		if err != nil {
			errCount++
			fmt.Printf("Error: HTTP failure, %s\n", err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errCount++
				fmt.Printf("Error: failed to read body of responce, %s\n", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					errCount++
					fmt.Printf("Status not 200 on /signup: status=%d, AT:%s\n", resp.StatusCode, godebug.LF())
				} else {
					if db1 || errCount > 0 {
						fmt.Printf("body ->%s<-\n", body)
					}
					err = json.Unmarshal(body, &tt)
					if err != nil {
						errCount++
						fmt.Printf("Error: failed to parse return JSON data, err=%s AT:%s data= ->%s<-\n", err, godebug.LF(), body)
					}
					if len(tt.Token) < 50 {
						errCount++
						fmt.Printf("Error: token is too short, AT:%s data= ->%s<-\n", godebug.LF(), body)
					}
				}
			}
		}
	}

	type User struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	type UserReturn struct {
		Users []User
	}

	// -----------------------------------------------------------------------------------------------------------------------
	// do Get on /users - should see our user in the list, with name
	{
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:18000/users", nil)
		req.Header.Add("User-Agent", "Test API Client")
		req.Header.Add("X-Authentication-Token", tt.Token)
		resp, err := client.Do(req)
		if err != nil {
			errCount++
			fmt.Printf("Error: HTTP failure, %s\n", err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errCount++
				fmt.Printf("Error: failed to read body of responce, %s\n", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					errCount++
					fmt.Printf("Status not 200 on /signup: status=%d, AT:%s\n", resp.StatusCode, godebug.LF())
				} else {
					if db1 || errCount > 0 {
						fmt.Printf("body ->%s<-\n", body)
					}
					var mdata UserReturn
					err = json.Unmarshal(body, &mdata)
					if err != nil {
						errCount++
						fmt.Printf("Error: failed to parse return JSON data, err=%s AT:%s data= ->%s<-\n", err, godebug.LF(), body)
					}
					// for user and original name
					if db1 || errCount > 0 {
						fmt.Printf("Parsed Data: %s\n", godebug.SVarI(mdata))
					}
					found := false
					for _, vv := range mdata.Users {
						if vv.Email == "xyz@test-example...." {
							found = true
							if vv.FirstName != "aa first" {
								errCount++
								fmt.Printf("Error: invalid name for user, AT:%s data= ->%s<-\n", godebug.LF(), body)
							}
						}
					}
					if !found {
						errCount++
						fmt.Printf("Error: failed to find the user, AT:%s data= ->%s<-\n", godebug.LF(), body)
					}

				}
			}
		}
	}

	// -----------------------------------------------------------------------------------------------------------------------
	// do PUT on /users - update name
	{
		client := &http.Client{}
		postData := map[string]string{"firstName": "bb first", "lastName": "bb last"}
		jsonValue, _ := json.Marshal(postData)
		req, err := http.NewRequest("PUT", "http://localhost:18000/users", bytes.NewReader(jsonValue))
		req.Header.Add("User-Agent", "Test API Client")
		req.Header.Add("Content-Type", "application/json") // set content type/encoding to JSON
		req.Header.Add("X-Authentication-Token", tt.Token)
		resp, err := client.Do(req)
		if err != nil {
			errCount++
			fmt.Printf("Error: HTTP failure, %s\n", err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errCount++
				fmt.Printf("Error: failed to read body of responce, %s\n", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					errCount++
					fmt.Printf("Status not 200 on /signup: status=%d, AT:%s\n", resp.StatusCode, godebug.LF())
				} else {
					if db1 || errCount > 0 {
						fmt.Printf("body ->%s<-\n", body)
					}
					err = json.Unmarshal(body, &tt)
					if err != nil {
						errCount++
						fmt.Printf("Error: failed to parse return JSON data, err=%s AT:%s data= ->%s<-\n", err, godebug.LF(), body)
					}
					if len(tt.Token) < 50 {
						errCount++
						fmt.Printf("Error: token is too short, AT:%s data= ->%s<-\n", godebug.LF(), body)
					}
				}
			}
		}
	}

	// -----------------------------------------------------------------------------------------------------------------------
	// do Get on /users - should see our user in the list, with new name
	{
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:18000/users", nil)
		req.Header.Add("User-Agent", "Test API Client")
		req.Header.Add("X-Authentication-Token", tt.Token)
		resp, err := client.Do(req)
		if err != nil {
			errCount++
			fmt.Printf("Error: HTTP failure, %s\n", err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errCount++
				fmt.Printf("Error: failed to read body of responce, %s\n", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					errCount++
					fmt.Printf("Status not 200 on /signup: status=%d, AT:%s\n", resp.StatusCode, godebug.LF())
				} else {
					if db1 || errCount > 0 {
						fmt.Printf("body ->%s<-\n", body)
					}
					var mdata UserReturn
					err = json.Unmarshal(body, &mdata)
					if err != nil {
						errCount++
						fmt.Printf("Error: failed to parse return JSON data, err=%s AT:%s data= ->%s<-\n", err, godebug.LF(), body)
					}
					// check for chaned name
					if db1 || errCount > 0 {
						fmt.Printf("Parsed Data: %s\n", godebug.SVarI(mdata))
					}
					found := false
					for _, vv := range mdata.Users {
						if vv.Email == "xyz@test-example...." {
							found = true
							if vv.FirstName != "bb first" {
								errCount++
								fmt.Printf("Error: invalid name for user, AT:%s data= ->%s<-\n", godebug.LF(), body)
							}
						}
					}
					if !found {
						errCount++
						fmt.Printf("Error: failed to find the user, AT:%s data= ->%s<-\n", godebug.LF(), body)
					}

				}
			}
		}
	}

	if errCount == 0 {
		fmt.Printf("PASS\n")
	} else {
		fmt.Printf("FAIL\n")
	}
}

const db1 = false

/* vim: set noai ts=4 sw=4: */
