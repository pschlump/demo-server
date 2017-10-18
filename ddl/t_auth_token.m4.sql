
-- drop TABLE "az_auth_token" ;
CREATE TABLE "az_auth_token" (
	"user_id"				uuid not null ,
	"auth_token"			uuid not null ,
	"created" 			timestamp default current_timestamp not null , 
	FOREIGN KEY ("user_id") REFERENCES "az_user" ("id")
);


-- Gurantee uniqueness on auth_tokens in this table.
CREATE UNIQUE INDEX "az_auth_token_u1" on "az_auth_token" ( "auth_token" );

-- User may have more than 1 auth token - multiple logins allowed.
CREATE UNIQUE INDEX "az_auth_token_u2" on "az_auth_token" ( "user_id", "auth_token" );

-- Allow Postgres 10x to do an index-only lookup of data in this table.  No table read occures.
CREATE UNIQUE INDEX "az_auth_token_p1" on "az_auth_token" ( "auth_token", "user_id" );

