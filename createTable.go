package main

//
// Demo server
//
// By Philip Schlump
// email: pschlump@gmail.com
// tel: 720-209-7888
//

import (
	"fmt"
	"os"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/demo-server/dbsql"
)

type CreateType struct {
	Name        string
	TableName   string // used for indexes
	Ddl         string
	DdlType     string
	IgnoreError bool
}

var createTable []CreateType

func init() {

	createTable = []CreateType{
		{ // 0
			Name:    "az_user",
			DdlType: "table",
			Ddl: `
CREATE TABLE "az_user" (
	"id"				uuid DEFAULT uuid_generate_v4() not null primary key,
	"email"				text not null,
	"password"			text not null,
	"first_name"		text not null,
	"last_name"			text not null,
	"created" 			timestamp default current_timestamp not null,
	"updated" 			timestamp
);
`,
		},
		{ // 1
			Name:      "az_user_u1",
			TableName: "az_user",
			DdlType:   "index",
			Ddl: `
CREATE UNIQUE INDEX "az_user_u1" on "az_user" ( "email" );
`,
		},
		{ // 2
			Name:    "az_user_upd",
			DdlType: "function",
			Ddl: `
CREATE OR REPLACE function az_user_upd()
RETURNS trigger AS 
$BODY$
BEGIN
  NEW.updated := current_timestamp;
  RETURN NEW;
END
$BODY$
LANGUAGE 'plpgsql';
`,
		},
		{ // 3
			Name:    "az_user_trig",
			DdlType: "trigger",
			Ddl: `
CREATE TRIGGER az_user_trig
BEFORE update ON "az_user"
FOR EACH ROW
EXECUTE PROCEDURE az_user_upd();
`,
		},
		{ // 4
			Name:    "az_auth_token",
			DdlType: "table",
			Ddl: `
CREATE TABLE "az_auth_token" (
	"user_id"				uuid not null ,
	"auth_token"			uuid not null ,
	"created" 			timestamp default current_timestamp not null , 
	FOREIGN KEY ("user_id") REFERENCES "az_user" ("id")
);
`,
		},
		{ // 5
			Name:      "az_auth_token_u1",
			TableName: "az_auth_token",
			DdlType:   "index",
			Ddl: `
CREATE UNIQUE INDEX "az_auth_token_u1" on "az_auth_token" ( "auth_token" );
`,
		},
		{ // 6
			Name:      "az_auth_token_u2",
			TableName: "az_auth_token",
			DdlType:   "index",
			Ddl: `
CREATE UNIQUE INDEX "az_auth_token_u2" on "az_auth_token" ( "user_id", "auth_token" );
`,
		},
		{ // 7
			Name:      "az_auth_token_p1",
			TableName: "az_auth_token",
			DdlType:   "index",
			Ddl: `
CREATE UNIQUE INDEX "az_auth_token_p1" on "az_auth_token" ( "auth_token", "user_id" );
`,
		},
		{ // 8
			Name:    "az_login",
			DdlType: "function",
			Ddl: `
CREATE or REPLACE FUNCTION az_login(p_email varchar, p_password varchar)
	RETURNS varchar AS $$
DECLARE
	   l_id 			uuid;
		l_auth_token	uuid;
	   l_data 			varchar (150);
	l_fail 			boolean;
BEGIN
	l_fail = false;
	l_id = null;
	l_data = '{ "status":"unknown"}';

	if not l_fail then
		select  "az_user"."id"
		into  l_id
			from "az_user" as "az_user" left join "az_auth_token" as "az_auth_token" on "az_auth_token"."user_id" = "az_user"."id"
			where "az_user"."email" = p_email
			  and "az_user"."password" = crypt(p_password, "az_user"."password")
			;

		if not found then
			l_data = '{ "status":"error", "code":"401", "msg":"Invalid username or password." }';
			l_fail = true;
		end if;
	end if;

	if not l_fail then
		l_auth_token = uuid_generate_v4();
		-- this is the palce to implement an expiration for the tokens - clean up of this table too.
		insert into "az_auth_token" (
			  "auth_token"
			, "user_id"
		) values (
			  l_auth_token
			, l_id
		);
		l_data = '{ "status":"success"'
			||', "auth_token":'||to_json(l_auth_token)
			||'}';
	end if;

	RETURN l_data;
END;
$$ LANGUAGE plpgsql;
`,
		},
		{ // 9
			Name:    "az_signup",
			DdlType: "function",
			Ddl: `
CREATE or REPLACE FUNCTION az_signup( p_email varchar, p_password varchar, p_first_name varchar, p_last_name varchar)
	RETURNS varchar AS $$
DECLARE
    l_id 				uuid;
	l_auth_token 		uuid;
	l_data				varchar (400);
	l_fail				bool;
BEGIN

	l_fail = false;
	l_data = '{"status":"success"}';
	l_id = uuid_generate_v4();

	IF not l_fail THEN
		BEGIN
			insert into "az_user" ( "id", "email", "password", "first_name", "last_name" )
				values ( l_id, p_email, crypt(p_password,gen_salt('bf',8)), p_first_name, p_last_name )
			;
		EXCEPTION
			WHEN unique_violation THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Unable to create user with this email address.  Please choose a different email address.","code":"601"}';
			WHEN others THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Database error occured.","code":"602"}';
		END;
	END IF;
	IF not l_fail THEN
		BEGIN
			IF not l_fail THEN
				l_auth_token = uuid_generate_v4();
				insert into "az_auth_token" (
					  "auth_token"	
					, "user_id"	
				) values (
					  l_auth_token
					, l_id
				);
			END IF;
		EXCEPTION
			WHEN unique_violation THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Unable to create user with this email address.  Please choose a different email address.","code":"603"}';
			WHEN others THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Database error occured.","code":"604"}';
		END;
	END IF;

	IF not l_fail THEN
		l_data = '{"status":"success"'
			||',"id":'||to_json(l_id)
			||',"auth_token":'||to_json(l_auth_token)
			||'}';
	END IF;

	RETURN l_data;
END;
$$ LANGUAGE plpgsql;
`,
		},
		{ // 10
			Name:    "az_update_user",
			DdlType: "function",
			Ddl: `
CREATE or REPLACE FUNCTION az_update_user( p_auth_token varchar,  p_first_name varchar, p_last_name varchar)
	RETURNS varchar AS $$
DECLARE
    l_id 				uuid;
	l_data				varchar (150);
	l_fail				bool;
BEGIN

	l_fail = false;
	l_data = '{"status":"unknown"}';

	IF not l_fail THEN
		BEGIN
			select "user_id"
				into l_id
				from "az_auth_token" 
				where "auth_token" = p_auth_token::uuid
			;
		IF not found THEN
			l_fail = true;
			l_data = '{ "status":"error", "code":"401", "msg":"Invalid auth_token - not found." }';
		END IF;
		EXCEPTION
			WHEN others THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Database error - invalid auth_token.","code":"602"}';
		END;
	END IF;

	IF not l_fail THEN
		BEGIN
			update "az_user"
				set "first_name" = p_first_name
				  , "last_name" = p_last_name
				where "id" = l_id
				;
		EXCEPTION
			WHEN others THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Database error occured.","code":"604"}';
		END;
	END IF;

	IF not l_fail THEN
		l_data = '{"status":"success"}';
	END IF;

	RETURN l_data;
END;
$$ LANGUAGE plpgsql;
`,
		},
	}

}

func ItemExistsDB(DdlType, Name, TableName string) bool {
	switch DdlType {
	case "table":
		TableName := Name
		qry := `SELECT * FROM information_schema.tables WHERE table_schema = $1 and table_name = $2`
		data, err := dbsql.SelData2(db, qry, "public", TableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError on table:%s, err=%s%s\n", MiscLib.ColorRed, TableName, err, MiscLib.ColorReset)
			return false
		} else if data == nil || len(data) == 0 {
			fmt.Fprintf(os.Stderr, "%sMissing table:%s%s\n", MiscLib.ColorYellow, TableName, MiscLib.ColorReset)
			return false
		}
		fmt.Fprintf(os.Stderr, "%sFound table: %s%s\n", MiscLib.ColorGreen, TableName, MiscLib.ColorReset)
		return true

	case "index":
		IndexName := Name
		qry := `
SELECT n.nspname as "Schema",
  c.relname as "Name",
  CASE c.relkind
		WHEN 'r' THEN 'table'
		WHEN 'v' THEN 'view'
		WHEN 'i' THEN 'index'
		WHEN 'S' THEN 'sequence'
		WHEN 's' THEN 'special'
	END as "Type",
  u.usename as "Owner",
 c2.relname as "Table"
FROM pg_catalog.pg_class c
     JOIN pg_catalog.pg_index i ON i.indexrelid = c.oid
     JOIN pg_catalog.pg_class c2 ON i.indrelid = c2.oid
     LEFT JOIN pg_catalog.pg_user u ON u.usesysid = c.relowner
     LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE c.relkind IN ('i','')
    AND n.nspname NOT IN ('pg_catalog', 'pg_toast')
    AND pg_catalog.pg_table_is_visible(c.oid)
	AND c.relkind = 'i'
	AND c2.relname = $3
	AND c.relname = $2
	AND '' <> $1
ORDER BY 1,2
			`
		data, err := dbsql.SelData2(db, qry, "public", IndexName, TableName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError on index:%s, err=%s%s\n", MiscLib.ColorRed, IndexName, err, MiscLib.ColorReset)
			return false
		} else if data == nil || len(data) == 0 {
			fmt.Fprintf(os.Stderr, "%sMissing index:%s%s\n", MiscLib.ColorYellow, IndexName, MiscLib.ColorReset)
			return false
		}
		fmt.Fprintf(os.Stderr, "%sFound index: %s%s\n", MiscLib.ColorGreen, IndexName, MiscLib.ColorReset)
		return true

	case "trigger":
		TriggerName := Name
		qry := `SELECT * FROM information_schema.triggers WHERE trigger_schema = $1 and trigger_name = $2`
		data, err := dbsql.SelData2(db, qry, "public", TriggerName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError on trigger:%s, err=%s%s\n", MiscLib.ColorRed, TriggerName, err, MiscLib.ColorReset)
			return false
		} else if data == nil || len(data) == 0 {
			fmt.Fprintf(os.Stderr, "%sMissing trigger:%s%s\n", MiscLib.ColorYellow, TriggerName, MiscLib.ColorReset)
			return false
		}
		fmt.Fprintf(os.Stderr, "%sFound trigger: %s%s\n", MiscLib.ColorGreen, TriggerName, MiscLib.ColorReset)
		return true

	case "function":
		FunctionName := Name
		qry := `SELECT routines.routine_name
			FROM information_schema.routines
			WHERE routines.specific_schema = $1
			  and ( routines.routine_name = lower($2)
			     or routines.routine_name = $2
				  )
			`
		data, err := dbsql.SelData2(db, qry, "public", FunctionName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError on function:%s, err=%s%s\n", MiscLib.ColorRed, FunctionName, err, MiscLib.ColorReset)
			return false
		} else if data == nil || len(data) == 0 {
			fmt.Fprintf(os.Stderr, "%sMissing function:%s%s\n", MiscLib.ColorYellow, FunctionName, MiscLib.ColorReset)
			return false
		}
		fmt.Fprintf(os.Stderr, "%sFound function: %s%s\n", MiscLib.ColorGreen, FunctionName, MiscLib.ColorReset)
		return true

	default:
		// Note this is missing views, types, sequences, constrints etc.
		// Also for functions should check parameters / return values
		// Also for tables should check columns
		// Also could have "alter" commands with IgnoreError==true to upgrade a schema.
		fmt.Fprintf(os.Stderr, "Invalid DdlType = ->%s<-\n", DdlType)
		os.Exit(5)
	}
	return false
}

// CreateTablesInDB creates tables and other objects in the database if they do not already exist.
// Returns error if required objects are missing.
func CreateTablesInDB() error {
	fail := false
	for ii, ddl := range createTable {
		if !ItemExistsDB(ddl.DdlType, ddl.Name, ddl.TableName) {
			err := dbsql.Run1(db, ddl.Ddl)
			if err != nil {
				if !ddl.IgnoreError {
					fail = true
					fmt.Fprintf(os.Stderr, "Error on %d item - creating %s of type %s, err=%s\n", ii, ddl.Name, ddl.DdlType, err)
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s/%s created\n", ddl.Name, ddl.DdlType)
			}
		}
	}
	if fail {
		return fmt.Errorf("Some objects were not found or created in the databse.")
	}
	return nil
}

/* vim: set noai ts=4 sw=4: */
