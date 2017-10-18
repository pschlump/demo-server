


-- DROP TABLE "az_user" ;
CREATE TABLE "az_user" (
	  "id"					uuid DEFAULT uuid_generate_v4() not null primary key
	, "email"				text not null
	, "password"		text not null
	, "first_name"			text not null
	, "last_name"			text not null
	, "created" 			timestamp default current_timestamp not null
	, "updated" 			timestamp
);

CREATE UNIQUE INDEX "az_user_u1" on "az_user" ( "email" );








CREATE OR REPLACE function az_user_upd()
RETURNS trigger AS 
$BODY$
BEGIN
  NEW.updated := current_timestamp;
  RETURN NEW;
END
$BODY$
LANGUAGE 'plpgsql';


CREATE TRIGGER az_user_trig
BEFORE update ON "az_user"
FOR EACH ROW
EXECUTE PROCEDURE az_user_upd();



