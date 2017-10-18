package main

//
// Demo server
//
// By Philip Schlump
// email: pschlump@gmail.com
// tel: 720-209-7888
//

var helpMsg = `
A simple server with /signon, /login, and /user endpoints.
This is a demo server that connects to Postgres and implements a simple API.
Once you have gotten the serve to connect to the database and run it with a
-C flag to create the database DDL you should be able to bring up the interactive
test pages.  The default is http://localhost:18000/index.html for the test pages.

Command Line Options:
    -p <Str> | --port=<Str>            The port to listen on or the IP:Port to listen on.
    -d <path> | --dir=<path>           Where to server static files from - for the demo page.
    -E | --errorstatus                 If specified errors are returned as a HTTP status
                                       like 401, not authorized.  Otherwise errors are 
                                       more detailed returned as a 200 with JSON of 
                                       { "status": "error", "msg": ..., "code": "..." }
    -D <str> | --debug=<str>           A comma separated list of debug flags.
    -C | --createddl                   If true then will create tables and other DDL.
    --log=<FileName>                   Where to send log output to.
    --cfg=<FileName>                   Input configuration file, ./cfg.jsonx default.
    --conn=<connect-str>               Connection string for Postgres, or default taken from 
                                       configuration file ./cfg.jsonx, -C|--cfg option.
    --help                             Print this message.

Debug Flags:

	To run with a flag set:
		demo-server -D db-func-call,IsLoggedIn
		demo-server -D echo-startup,echo-cfg,db-func-call

	Flags are:
		CallFunction				Output information on call to stored procedure.
		CheckAuthTokenInDb			Check of auth_token in database
		CreateJwtToken				Creation of JWT token
		IsLoggedIn					Check of login status.
		db-func-call				Print detailed information on stored procedure call.
		db-jwt-token				Print details of JWT token and signature.
		db-token					Print out the auth_token information.
		db-validate-token			Print out validation of auth token/JWT token.
		dump-meta					Print out meta-data from the JsonX call.
		echo-cfg					Print out configuration file.
		echo-startup				Print out configuration before start of server.

Configuration File: ./cfg.jsonx is in a superset of JSON that allows for substitution
of environment variables with {{ __env__ NAME }}.  This is used to allow you to keep
you Postgres password in your environment instead of putting it in a file that might
get checked into github.com.   For a full set of documentation on the configuration
system see: https://github.com/pschlump/JSONx.   The Fields in the  configuration file
are:

    comment          A comment to help identify the file
    PGConn           The Postgres connection string.  Note the """ quotes.  Can be
                     overridden with the --conn command line option.
    DBType           Database type should be: "postgres".
    DBName           Database name - not used by Postgres. Default "pschlump".
    KeyFile          The public key file, defaults to ./key/sample_key.pub.
    KeyFilePrivate   The private key file, defaults to ./key/sample_key.

`

/* vim: set noai ts=4 sw=4: */
