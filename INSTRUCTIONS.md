Instructions for Installing and Running
===================================================================

## Install the Software

You have probably already done the "clone" line.

```bash
$ git clone https://github.com/pschlump/demo-server.git
$ cd demo-server
$ go get
```

The software should be complete.  There is a web page for interactive testing.
It includes all the style sheets and JavaScript to run.

## Start Testing of Code 

In the `demo-server` directory

```bash
$ cd misc
$ go get
$ go test
```

You should run and see a *PASS*.

Now change to `../jwtlib`.

```bash
$ cd ../jwtlib
$ go get
$ go test
```

You should run and see a *PASS*.

## Configure the Database.

The code uses UUIDs and the "crypto" extension in Postgres.  This will require that
these extensions are available.

In the `./demo-server/ddl` directory there is a short script create the
two extensions.  Postgres extensions have to be installed on the system.  With a default
install of Postgres on Ubuntu you can `sudo apt-get install postgresql-contrib-9.6`
to install the extensions.  You will need to change 9.6 to the correct version 
for your database.

```bash
$ cd ./ddl
$ psql -p 5432 -a -P pager=off -h {{IP-of-Pg-Server}} -U {{your-username}} --dbname={{your-database}}
username=# \i setup.sql
username=# \q
$ cd .. 
```

## Configure to Talk to the Database.

In the `demo-server` directory there should be a `cfg.jsonx` file.
Edit it.  You will need to change the line that reads:

```
	PGConn: """user=pschlump password={{ __env__ DB_PASS }} sslmode=disable dbname=pschlump port=5432 host=192.168.0.139""" ,
```

to have a correct Postgres connection string in it.  If you keep your database password in 
your environment, then change the `{{ __env__ DB_PASS }}` to use the correct environment variable
for the database password.  You can also hard code the password into the connection string.
This is dangerous.  If a non-authorized person gains access to the file, say by checking
the file into github.com, then all sorts of bad things could happen.

Now change to `../dbsql` directory.

```bash
$ cd ../dbsql
$ go get
$ go test
```

The test should pass.  A `PASS` will indicate that you have connected to the database and
that the database interface library works.

I think that this is a big milestone.  Give yourself a pat on the back.

## Create DDL in the Database (Tables, etc...)

The server, `demo-server` has a command line option to create all of the
tables that are needed in the database.  This is the -C option.

```bash
$ cd ..
$ pwd
```

You should be in the `.../demo-server` directory.


```
$ go get
$ go build
```

You should have a compiled executable named `demo-server`.

Run it with the extra option.

```
$ ./demo-server -C
```

It should print out some colored output with yellow lines indicating that it did not find
some database objects and green lines indicating that it created the objects.  If it fails to
create all the objects then it will exit.   If you see:

```
...
...
...
Successfully connected to Postgres and Listining on 18000
```

On future runs of the server you will not need the -C option.

## Test the API

then the server is running on port 18000.   There are command line options to change the
port and IP address that the server listens to.  Use the `--help` option if you need to 
pick a different port.

In a new window change into the `demo-server/testApi` directory.
This is an automated test program for the API.  Let's build it and run it.

```
$ go get
$ go build
$ ./testApi
```

This test assumes that the server is running on `localhost:18000`.  If you have 
changed the IP or the port, then you will need to edit the source to change the
test.

If the API is working, you should see a *PASS*.

You can also interactively test the API from a web browser.  Use the URL
`http://localhost:18000/` for the interactive test.

## Diagnosing Problems

Did everything compile?  If not, then call me and tell me what failed.

Did all of the stored procedures, tables, indexes and other DDL get created?
If not, the source is in `./ddl` for all of them.  Connect to Postgres with `psql`
and manually create the missing items.
You should have 3 tables:  `az_output`, `az_user`, and `az_auth_token`.
You should have 4 functions: `az_update_user`, `az_signup` `az_login`, and `az_user_upd`.
You should have 1 trigger: `az_user_trig` and a bunch of indexes.

Run the automated test for the stored procedures.  This is the file `test.sql`
in `./ddl`.  Run in `psql`.   It should produce *PASS* and then some output
as to the set of test that it ran.   After running this
you should have 3 tables:  `az_output`, `az_user`, and `az_auth_token`.

## Notes

The `./ddl` directory has the source for all the tables/indexes and other database
objects.

The `./key` directory has a pair of self signed RSA 256 bit keys that should only
be used for testing.    These are the same keys that are distributed with the JWT go
library.  You can create new keys if you want.

The static files for the web test are in `./www`.

The Logs from running the server are in `./log`.

## Contact Information

Philip Schlump
1005 E. Custer St. 
Laramie WY, 82070
tel: 720-209-7888 (cell)
email: pschlump@gmail.com

