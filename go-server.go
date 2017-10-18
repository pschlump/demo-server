//
// Demo server
//
// By Philip Schlump
// email: pschlump@gmail.com
// tel: 720-209-7888
//

package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	_ "github.com/lib/pq"

	JsonX "github.com/pschlump/JSONx"
	"github.com/pschlump/demo-server/dbsql"
	"github.com/pschlump/demo-server/jwtlib"
	"github.com/pschlump/demo-server/misc"
	"github.com/pschlump/godebug"
)

// Note: https://github.com/dgrijalva/jwt-go.git // jwt "github.com/dgrijalva/jwt-go"

var NameVersion = "demo-server version 0.2.1"
var NameProgram = "demo-server"

var Port = flag.String("port", "", "Port to listen on, or ip:port")                                             // 1
var Debug = flag.String("debug", "", "Debug flags")                                                             // 3
var Dir = flag.String("dir", "", "Directory to server static files from, default ./www")                        // 8
var CreateDDL = flag.Bool("createddl", false, "Create Table DDL")                                               // 7
var ErrorStatus = flag.Bool("errorstatus", false, "Return errors as HTTP code, --errorstatus, or JSON default") // 14
var Help = flag.Bool("help", false, "Print help message")                                                       // 15
var LogFile = flag.String("log", "", "log errors to this file")                                                 // 9

var CfgFile = flag.String("cfg", "cfg.jsonx", "config file")       // 10
var PGConn = flag.String("conn", "", "PotgresSQL connection info") // 11 Xyzzy -- test from CLI

func init() {
	flag.StringVar(Port, "p", "", "Port to listen on, or ip:port")                           // 1
	flag.StringVar(Debug, "D", "", "Debug flags")                                            // 3
	flag.StringVar(Dir, "d", "", "Directory to server static files from, default ./www")     // 8
	flag.BoolVar(CreateDDL, "C", false, "Create Table DDL")                                  // 7
	flag.BoolVar(ErrorStatus, "E", false, "Return errors as HTTP code, -E, or JSON default") // 14
	flag.StringVar(LogFile, "l", "", "log errors to this file")                              // 9
}

var db *sql.DB
var dbName = "" // not used for Postgres -- Important for mySql/MariaDB, microsoft SQL server
var ferr *os.File
var fo *os.File

type CfgType struct {
	Comment        string // not used - the config file has comment in it saying what it is and what program it belongs to
	Dir            string `gfJsonX:"Dir" gfDefault:"./www"`
	Port           string `gfJsonX:"Port" gfDefault:"18000"`
	PGConn         string `gfJsonX:"PGConn"`
	DBType         string `gfJsonX:"DBType" gfDefault:"postgres"`
	DBName         string `gfJsonX:"DBName" gfDefault:"pschlump"`
	KeyFile        string `gfJsonX:"KeyFile" gfDefault:"./key/sample_key.pub"`
	KeyFilePrivate string `gfJsonX:"KeyFile" gfDefault:"./key/sample_key"`
}

var Cfg CfgType

// -------------------------------------------------------------------------------------------------
func main() {

	var err error

	fo = os.Stdout
	ferr = os.Stderr

	flag.Parse()
	fns := flag.Args()
	if len(fns) != 0 || *Help {
		fmt.Fprintf(os.Stderr, "Usage: %s [...options...]\n%s", NameProgram, helpMsg)
		os.Exit(1)
	}

	JsonX.SetDebugFlags(*Debug)
	misc.SetDebugFlags(*Debug)

	if *LogFile != "" {
		ferr, err = misc.Fopen(*LogFile, "a")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to open %s for outout, %s\n", *LogFile, err)
			os.Exit(3)
		}
		fmt.Fprintf(ferr, `{"startup-go-server":%q}`+"\n", time.Now().Format(time.RFC3339Nano))
	}

	auth := ""

	if !misc.Exists(*CfgFile) {
		fmt.Fprintf(os.Stderr, "%s: Configuration file is required, file name=%s\n", NameProgram, *CfgFile)
		os.Exit(1)
	} else {

		fn := *CfgFile
		In, err := ioutil.ReadFile(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error returned from ReadFile: %s\n", NameProgram, err)
			os.Exit(1)
		}

		meta, err := JsonX.Unmarshal(fn, In, &Cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error returned from JsonX.Unmarshal: %s\n", NameProgram, err)
			os.Exit(1)
		}

		if misc.Db["dump-meta"] {
			// _ = meta // Throw away meta info on conversion of JsonX to data structure.
			fmt.Printf("AT: %s, meta = %s\n", godebug.LF(), misc.SVarI(meta))
		}

		if misc.Db["echo-cfg"] {
			fmt.Printf("Read in config file, cfg=%s\n", misc.SVarI(Cfg))
		}

		if *PGConn != "" {
			auth = *PGConn
		} else if Cfg.PGConn != "" {
			auth = Cfg.PGConn
		}
		Cfg.PGConn = "...removed for security reasons..."

		if *Dir != "" {
			// fmt.Printf("*Dir=%s At:%s\n", *Dir, godebug.LF())
		} else if Cfg.Dir != "" {
			*Dir = Cfg.Dir
			// fmt.Printf("*Dir=%s At:%s\n", *Dir, godebug.LF())
		}

	}

	fmt.Fprintf(ferr, `{"config":%s}`+"\n", misc.SVarI(Cfg))

	if !misc.Exists(Cfg.KeyFile) {
		fmt.Fprintf(os.Stderr, "%s: Public key pair is required for validation of JWT token signature.  %s should be public key - not found.\n", NameProgram, Cfg.KeyFile)
		os.Exit(1)
	}
	if !misc.Exists(Cfg.KeyFilePrivate) {
		fmt.Fprintf(os.Stderr, "%s: Private key pair is required for validation of JWT token signature.  %s should be public key - not found.\n", NameProgram, Cfg.KeyFilePrivate)
		os.Exit(1)
	}

	if auth == "" {
		fmt.Fprintf(os.Stderr, `{"error":"Connection information to postgres is required."}`+"\n")
		os.Exit(1)
	} else {
		db_x := dbsql.ConnectToAnyDb("postgres", auth, dbName)
		if db_x == nil {
			fmt.Fprintf(os.Stderr, "%s: Unable to connection to database: %v\n", NameProgram, err)
			os.Exit(1)
		}
		db = db_x.Db
	}

	// check on the existence and create tables as necessary in the database.
	if *CreateDDL {
		fail := CreateTablesInDB()
		if fail != nil {
			os.Exit(1)
		}
	}

	*Dir = path.Clean(*Dir)

	if *Dir == "" || *Dir == "." || *Dir == "./" {
		TDir, _ := os.Getwd()
		*Dir = TDir
	}
	if *Port == "" {
		*Port = "18000"
	}

	if misc.Db["echo-startup"] {
		fmt.Printf("cfg=%s\n", misc.SVarI(Cfg))
		fmt.Printf("serving files from %s\n", *Dir)
		fmt.Printf("connected to Postgres\n")
	}

	http.HandleFunc("/api/status", respHandlerStatus)

	// POST /signup - create new user (register) - return jwt - login
	// POST /login - login - return jwt
	// GET /users	- list of users
	// PUT /users	- update users - based on jwt token

	http.HandleFunc("/signup", handleSignup)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/users", handleUsers)
	http.Handle("/", http.FileServer(http.Dir(*Dir)))

	fmt.Fprintf(os.Stderr, "Successfully connected to Postgres and Listining on %s\n", *Port)

	log.Fatal(http.ListenAndServe(":"+(*Port), nil))
}

// -------------------------------------------------------------------------------------------------

type JwtToken struct {
	Token string
}

type Login struct {
	Email    string
	Password string
}

type User struct {
	Email     string
	Password  string
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserList struct {
	Users []User
}

type Register struct {
	Email     string
	Password  string
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// IsLoggedIn returns true if the `X-Authentication-Token` is a valid JWT signed token and
// if the JWT token's claims include a valid auth_token.
func IsLoggedIn(req *http.Request) (auth_token string, ok bool) {
	ok = true
	qq := req.RequestURI
	jwt_token := req.Header.Get(`X-Authentication-Token`)
	if misc.Db["IsLoggedIn"] {
		fmt.Printf("jwt_token = [%s], AT: %s\n", jwt_token, godebug.LF())
	}
	authToken, err := jwtlib.VerifyToken([]byte(jwt_token), Cfg.KeyFile)
	if err != nil {
		ok = false
	}
	if misc.Db["IsLoggedIn"] {
		fmt.Printf("Before CheckAuthTokenInDb, ok = %v authToken = [%s], AT: %s\n", ok, authToken, godebug.LF())
	}
	if ok {
		ok = CheckAuthTokenInDb(authToken) // if ( not a valid auth_token in database )
	}
	if misc.Db["IsLoggedIn"] {
		fmt.Printf("ok=%v authToken=[%s], AT: %s\n", ok, authToken, godebug.LF())
	}
	if !ok {
		if misc.Db["IsLoggedIn"] {
			fmt.Printf("jwt_token = [%s], AT: %s\n", jwt_token, godebug.LF())
		}
		if req.Method == "GET" {
			fmt.Fprintf(ferr, `{"error":"login","msg":"attempt to access %s:%s without login","reqHeaders":%s}`+"\n", req.Method, qq, misc.SVar(req.Header))
			// Xyzzy - get all URL params
		} else {
			fmt.Fprintf(ferr, `{"error":"login","msg":"attempt to access %s:%s without login","reqHeaders":%s}`+"\n", req.Method, qq, misc.SVar(req.Header))
			// Xyzzy - dump out body also
		}
	}
	return authToken, ok
}

// CheckAuthTokenInDb looks up the authToken in the database and determines if it is still a valid
// token.  A token could be deleted after it has been issued.   This takes care of that case.
// True is returned if valid.
func CheckAuthTokenInDb(authToken string) bool {
	stmt := `select "user_id" from "az_auth_token" where "auth_token" = $1::uuid`
	if misc.Db["CheckAuthTokenInDb"] {
		fmt.Printf("authToken=%s AT: %s\n", authToken, godebug.LF())
	}
	rows, err := dbsql.SelData2(db, stmt, authToken)
	if err != nil {
		fmt.Fprintf(ferr, `{"error":"database","msg":%q,"stmt":%q,"auth_token":%q}`+"\n", "not a valid token", stmt, authToken)
		return false
	}
	if len(rows) != 1 {
		fmt.Fprintf(ferr, `{"error":"database","msg":%q,"stmt":%q,"auth_token":%q}`+"\n", "invalid number of rows", stmt, authToken)
		return false
	}
	return true
}

// handleSignup creates a new user in the database and returns a token if succesful.
// Input body:
//
// 	{
// 	  "email": "test@example.com",
// 	  "password": "aPassword",
// 	  "firstName": "FirstName",
// 	  "lastName": "LastName"
// 	}
//
// where `email` is a unique key in the database.
//
// The response body should return a JWT on success that can be used for other endpoints:
//
// 	{
// 	  "token": "some_jwt_token"
// 	}
//
func handleSignup(www http.ResponseWriter, req *http.Request) {
	var rv string

	switch req.Method {
	case "POST":

		var data Register
		LogRequest(www, req)

		if !ParseBody(www, req, &data) {
			return
		}

		if data.Email == "" {
			ErrorReturn(0, www, req, http.StatusBadRequest, "BadData", "Email Required", "Email address can not be empty")
			return
		}
		if data.Password == "" {
			ErrorReturn(0, www, req, http.StatusBadRequest, "BadData", "Password Required", "Password address can not be empty")
			return
		}
		if data.FirstName == "" {
			ErrorReturn(0, www, req, http.StatusBadRequest, "BadData", "First Name Required", "FirstName address can not be empty")
			return
		}
		if data.LastName == "" {
			ErrorReturn(0, www, req, http.StatusBadRequest, "BadData", "Last Name Required", "LastName address can not be empty")
			return
		}

		stmt := `select az_signup($1,$2,$3,$4) as "x"` // email, password, firs_name, last_name -> auth_token
		funcRv, err := CallFunction(stmt, data.Email, data.Password, data.FirstName, data.LastName)
		// fmt.Printf("%sfuncRv = %s, AT:%s%s\n", MiscLib.ColorRed, godebug.SVarI(funcRv), godebug.LF(), MiscLib.ColorReset)
		if err != nil {
			ErrorReturn(0, www, req, http.StatusBadRequest, "FuncCall", "Database Error", "Function call failed, stmt=%s, data=%s", stmt, misc.SVar(data))
			// fmt.Fprintf(ferr, `{"error":"Call Function","msg":%q,"stmt":%q,"data":%s}`+"\n", err, stmt, misc.SVar(data))
			return
		}
		if funcRv.Status != "success" {
			ErrorReturn(0, www, req, http.StatusBadRequest, "FuncCall", funcRv.Msg, "Function call failed, stmt=%s, data=%s funcRv=%s", stmt, misc.SVar(data), misc.SVar(funcRv))
			return
		}
		rv, err = jwtlib.CreateJwtToken(funcRv.AuthToken, Cfg.KeyFilePrivate)
		if err != nil {
			ErrorReturn(0, www, req, http.StatusBadRequest, "JwtTokenErr", "Toeken Error", "Unable to create JWT token, err=%q", err)
			return
		}

	default:
		ErrorReturn(0, www, req, http.StatusBadRequest, "ReqMethod", "Invalid Method", "Invalid request method")
		return
	}

	//	www.WriteHeader(200)
	io.WriteString(www, rv)
}

// handleLogin hancles POST:/login returning a JWT token if login is successful.
// Input body:
//
//	{
//	  "email": "test@example.com",
//	  "password": "aPassword"
//	}
//
// The response body should return a JWT on success that can be used for other endpoints:
//
//	{
//	  "token": "some_jwt_token"
//	}
//
func handleLogin(www http.ResponseWriter, req *http.Request) {
	var rv string

	switch req.Method {
	case "POST":
		var data Login
		LogRequest(www, req)

		if !ParseBody(www, req, &data) {
			return
		}

		if data.Email == "" {
			ErrorReturn(0, www, req, http.StatusBadRequest, "BadData", "Email Required", "Email address can not be empty")
			return
		}
		if data.Password == "" {
			ErrorReturn(0, www, req, http.StatusBadRequest, "BadData", "Password Required", "Password address can not be empty")
			return
		}

		stmt := `select az_login($1,$2) as "x"` // email, password
		funcRv, err := CallFunction(stmt, data.Email, data.Password)
		if err != nil {
			ErrorReturn(0, www, req, http.StatusBadRequest, "FuncCall", "Database Error", "Function call failed, stmt=%s, data=%s", stmt, misc.SVar(data))
			return
		}
		if funcRv.Status != "success" {
			ErrorReturn(0, www, req, http.StatusBadRequest, "FuncCall", funcRv.Msg, "Function call failed, stmt=%s, data=%s funcRv=%s", stmt, misc.SVar(data), misc.SVar(funcRv))
			return
		}
		rv, err = jwtlib.CreateJwtToken(funcRv.AuthToken, Cfg.KeyFilePrivate)
		if err != nil {
			ErrorReturn(0, www, req, http.StatusBadRequest, "JwtTokenErr", "Toeken Error", "Unable to create JWT token, err=%q", err)
			return
		}

	default:
		ErrorReturn(0, www, req, http.StatusBadRequest, "ReqMethod", "Invalid Method", "Invalid request method")
		return
	}

	//	www.WriteHeader(200) // validate that this is an success
	io.WriteString(www, rv)
}

// handleUsers responds to GET:/users and PUT:/users.  GET returns a list of users.
// PUT updates the first/last name fields for users.
//
// GET:/users
//
// Endpoint to retrieve a json of all users. This endpoint requires a valid `x-authentication-token` header to be passed in
// with the request.
//
// The response body should look like:
//
// 	{
// 	  "users": [
// 		{
// 		  "email": "test@example.com",
// 		  "firstName": "FirstName",
// 		  "lastName": "LastName"
// 		}
// 	  ]
// 	}
//
//
// PUT:/users
// Endpoint to update the current user `firstName` or `lastName` only. This endpoint requires a valid
// `x-authentication-token` header to be passed in and it should only update the user of the JWT being passed in. The
// payload can have the following fields:
//
//
// 	{
// 	  "firstName": "NewFirstName",
// 	  "lastName": "NewLastName"
// 	}
//
func handleUsers(www http.ResponseWriter, req *http.Request) {
	var rv string

	switch req.Method {
	case "GET":
		LogRequest(www, req)

		if _, ok := IsLoggedIn(req); !ok { // validate user is logged in - check for header x-auth... validate JWT token.

			ErrorReturn(0, www, req, http.StatusUnauthorized, "Unauth", "Must login first", "Login Requried for GET:/users")

		} else {

			stmt := `select "email", "first_name" as "firstName", "last_name" as "lastName" from "az_user" order by 1`
			rows, err := dbsql.SelData2(db, stmt)
			if err != nil {
				rv = fmt.Sprintf(`{"status":"error","msg":"Unable to get the set of users.   Database error has been logged.\n"}`)
				fmt.Fprintf(ferr, `{"error":"database","msg":%q,"stmt":%q,"reqHeaders":%s}`+"\n", err, stmt, misc.SVar(req.Header))
			} else {
				type ReturnData struct {
					Users []map[string]interface{} `json:"Users"`
				}
				var rd ReturnData
				rd.Users = rows
				rv = misc.SVarI(rd)
			}
		}

	case "PUT":
		LogRequest(www, req)

		if authToken, ok := IsLoggedIn(req); !ok { // validate user is logged in - check for header x-auth... validate JWT token.

			ErrorReturn(0, www, req, http.StatusUnauthorized, "Unauth", "Must login first", "Login Requried for PUT:/users")

		} else {
			var data UpdateUser
			LogRequest(www, req)

			if !ParseBody(www, req, &data) {
				return
			}

			if misc.Db["IsLoggedIn"] {
				fmt.Printf("In Update: authToken=[%s], AT: %s\n", authToken, godebug.LF())
			}

			if data.FirstName == "" {
				ErrorReturn(0, www, req, http.StatusBadRequest, "BadData", "First Name Required", "FirstName address can not be empty")
				return
			}
			if data.LastName == "" {
				ErrorReturn(0, www, req, http.StatusBadRequest, "BadData", "Last Name Required", "LastName address can not be empty")
				return
			}

			// CREATE or REPLACE FUNCTION az_update_user( p_auth_token varchar,  p_first_name varchar, p_last_name varchar)
			stmt := `select az_update_user ( $1, $2, $3 ) as "x"` // auth_token, first_name, last_name  -- join to auth_token, get count of rows updated
			funcRv, err := CallFunction(stmt, authToken, data.FirstName, data.LastName)
			if err != nil {
				ErrorReturn(0, www, req, http.StatusBadRequest, "FuncCall", "Database Error", "Error on call to function, stmt=%s, data=%s, err=%s", stmt, misc.SVar(data), err)
				return
			}

			// xyzzy - this way
			if funcRv.Status == "success" {
				rv = `{"status":"success"}`
			} else {
				ErrorReturn(0, www, req, http.StatusBadRequest, "UpdateErr", "Database Error", "Error on call to function, data=%s", misc.SVar(funcRv))
				return
			}

		}

	default:
		ErrorReturn(0, www, req, http.StatusBadRequest, "ReqMethod", "Invalid Method", "Invalid request method")
		return
	}

	//	www.WriteHeader(200) // validate that this is an success
	io.WriteString(www, rv)
}

// respHandlerStatus is an "echo" of what was sent from the client side.  This can be used to
// check what is transfered from the client in a request or as a "ping" type response to verify
// that the server is alive.
func respHandlerStatus(res http.ResponseWriter, req *http.Request) {
	qq := req.RequestURI

	var rv string
	res.Header().Set("Content-Type", "application/json")
	rv = fmt.Sprintf(`{"status":"success","name":%q,"URI":%q,"req":%s, "response_header":%s}`, NameVersion, qq, misc.SVarI(req), misc.SVarI(res.Header()))

	io.WriteString(res, rv)
}

// Data returned from the function calls.
type FuncRv struct {
	Status    string // success, error, unkown, fail
	Msg       string // if not Success then message
	Code      string // numeric code for errors - helps with internationalization
	AuthToken string `json:"auth_token"`
	// Id        string `json:"id"` // user_id - discarded
}

// CallFunction calls a database stored procedure.  The stored procedure must return a JSON string.
// The stmt must be of the form:
//    select function_name(params) as "x"
// Bind variables can be passed as `args`.
// The return value of the function is unmarshaled using FuncRv.
func CallFunction(stmt string, args ...interface{}) (rv FuncRv, err error) {
	rows, err := dbsql.SelData2(db, stmt, args...)
	if err != nil {
		rv = FuncRv{Status: "error", Msg: "Unable to the set of users.   Database error logged.\n", Code: "100"}
		fmt.Fprintf(ferr, `{"error":"database","msg":%q,"stmt":%q}`+"\n", err, stmt)
	} else {
		// rows: [{"x":"{\"status\":\"success\",\"id\":\"c5714cb0-5a25-4dd8-a327-7957a7d12f8e\",\"auth_token\":\"eb840017-b6ad-45a1-aafb-6ceb84a29e91\"}"}]
		if misc.Db["CallFunction"] {
			fmt.Printf("Rows: %s, %s\n", misc.SVar(rows), godebug.LF())
		}
		if len(rows) == 0 {
			err = fmt.Errorf("Function call did not return a row") // Xyzzy - log
			return
		}
		if len(rows) > 1 {
			err = fmt.Errorf("Function call returned too many rows, %d expected 1", len(rows)) // Xyzzy - log
			return
		}
		xx0, ok := rows[0]["x"]
		if !ok {
			err = fmt.Errorf("Function call did not have a column called \"x\", the return value") // Xyzzy - log
			return
		}
		xx, ok := xx0.(string)
		if !ok {
			err = fmt.Errorf("Function call can not convert return value from function to a string") // Xyzzy - log
			return
		}
		err = json.Unmarshal([]byte(xx), &rv)
		if err != nil {
			err = fmt.Errorf("Unable to JSON parse the return value from %s, err=%s, in %s\n", stmt, err, xx) // Xyzzy - log
		}
		if misc.Db["db-func-call"] {
			fmt.Printf("stmt=%s raw=%s funcRv=%s args=%s\n", stmt, xx, misc.SVar(rv), args)
		}
	}
	return
}

// ParseBody parses JSON bodies.  This assuems that all bodies are JSON *not* URL/Form encoded.
func ParseBody(www http.ResponseWriter, req *http.Request, data interface{}) (ok bool) {
	ok = true
	if req.Body == nil {
		ErrorReturn(1, www, req, http.StatusBadRequest, "BadBody", "Requst must have a body", "body is nil")
		return false
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ErrorReturn(1, www, req, http.StatusBadRequest, "BadBody", "Invalid body", "body is unreadable, err=%s", err)
		return false
	}
	defer req.Body.Close()
	err = json.Unmarshal(body, data)
	if err != nil {
		ErrorReturn(1, www, req, http.StatusBadRequest, "ParseErr", "Invalid body", "Parse error in JSON body, err=%s body=%s", err, body)
		return false
	}
	fmt.Fprintf(ferr, `{ "body":%q }`+"\n", body)
	return
}

// LogRequest logs the reuest to the log file.
func LogRequest(www http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(ferr, `{ "method":%q, "URI":%s, "AT":%q, "time":%q }`+"\n", req.Method, req.RequestURI, godebug.LF(2), time.Now().Format(time.RFC3339Nano))
}

// ErrorReturn logs errors and returns either a status code or a 200/message
func ErrorReturn(n int, www http.ResponseWriter, req *http.Request, status int, errType, umsg, msg string, args ...interface{}) {
	var rv string
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	fmt.Fprintf(ferr, `{ "error":%q, "status":%d, "msg":%q, "AT":%q, "method":%q, "uri":%q, "time":%q }`+"\n", errType, status, godebug.LF(2+n), msg, req.Method, req.RequestURI, time.Now().Format(time.RFC3339Nano))
	if *ErrorStatus {
		www.WriteHeader(status)
		return
	}
	rv = fmt.Sprintf(`{"status":"error","msg":%q,"statusCode":%d}`, umsg, status)
	www.WriteHeader(http.StatusOK)
	io.WriteString(www, rv)
}

/* vim: set noai ts=4 sw=4: */
