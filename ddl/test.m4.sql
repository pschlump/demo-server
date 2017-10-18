
drop FUNCTION az_test_01();
-- drop TABLE "az_output" ;

CREATE SEQUENCE az_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;

CREATE TABLE "az_output" (
	  "seq"	 		bigint DEFAULT nextval('az_id_seq'::regclass) NOT NULL 
	, "msg"			text
);


CREATE or REPLACE FUNCTION az_test_01()
	RETURNS varchar AS $$
DECLARE
    l_id 			uuid;
  	l_auth_token	uuid;
  	l_user_id		uuid;
    l_data 			varchar (150);
    l_flag 			varchar (150);
    l_rv 			varchar (1000);
	l_fail 			boolean;
	l_password 		varchar(10);
	l_n1 			int;
	l_n2 			int;
    l_old_name 		varchar (150);
    l_new_name 		varchar (150);
    l_email 		varchar (150);
BEGIN

	l_fail = false;
	l_data = 'PASS';
	l_password = 'aa';

	insert into "az_output" ( "msg" ) values ( 'Start of Test ---------------------------------------------------------------------------------' );

	BEGIN
		l_id = uuid_generate_v4();
		insert into "az_output" ( "msg" ) values ( '  PASS - genrate uuid' );
 	EXCEPTION
 		WHEN others THEN
 			l_fail = true;
 			l_data = '{"status":"error","msg":"Unable to generate UUIDs - may need to run ./ddl/setup.sql"}';
	END;

	insert into "az_output" ( "msg" ) values ( '  l_id='||(l_id::varchar) );

	l_email = l_id::varchar || '@test-example....';

	IF not l_fail THEN
		BEGIN
			delete 
				from "az_user"
				where "email" like '%@test-example....'
				;
			insert into "az_output" ( "msg" ) values ( '  PASS - cleanup of old test user.' );
			select count(1)
				into l_n1
				from "az_auth_token"
				;
			select az_signup( l_email, l_password, 'aa first', 'aa last')
				into l_rv
				;
			insert into "az_output" ( "msg" ) values ( '  RAN - az_signup' );
			insert into "az_output" ( "msg" ) values ( '    l_rv='||l_rv );
			select 'found' as "x"
				into l_flag
				from "az_user"
				where "email" like '%@test-example....'
				;
			if not found then
				l_fail = true;
				l_data = '{"status":"error","msg":"az_signup failed to insert a user - case 1"}';
			end if;
			if l_flag <> 'found' then
				l_fail = true;
				l_data = '{"status":"error","msg":"az_signup failed to insert a user - case 2"}';
			end if;
			select count(1)
				into l_n2
				from "az_auth_token"
				;
			if l_n1 >= l_n2 then
				l_fail = true;
				l_data = '{"status":"error","msg":"az_login failed to insert create an auth token."}';
			end if;
		EXCEPTION
			WHEN others THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Error occured on az_signup"}';
		END;
	END IF;

	IF not l_fail THEN
		BEGIN
			select count(1)
				into l_n1
				from "az_auth_token"
				;
			select az_login( l_email, l_password )
				into l_rv
				;
			insert into "az_output" ( "msg" ) values ( '  RAN - az_login' );
			insert into "az_output" ( "msg" ) values ( '    l_rv='||l_rv );
			select count(1)
				into l_n2
				from "az_auth_token"
				;
			if l_n1 >= l_n2 then
				l_fail = true;
				l_data = '{"status":"error","msg":"az_login failed to insert create an auth token."}';
			end if;
		EXCEPTION
			WHEN others THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Error occured on az_login"}';
		END;
	END IF;

	IF not l_fail THEN
		BEGIN
			select "first_name", "id"
				into l_old_name, l_user_id
				from "az_user"
				where "email" = l_email
				;
			select "auth_token"
				into l_auth_token
				from "az_auth_token"
				where "user_id" = l_user_id
				;
			select az_update_user( l_auth_token::varchar, 'Timmothy', 'Smith' )
				into l_rv
				;
			insert into "az_output" ( "msg" ) values ( '  RAN - az_update_user' );
			insert into "az_output" ( "msg" ) values ( '    l_rv='||l_rv );
			select "first_name"
				into l_new_name
				from "az_user"
				where "email" = l_email
				;
			if l_old_name <> 'aa first'  then
				l_fail = true;
				l_data = '{"status":"error","msg":"az_update_user incorect first name."}';
			end if;
			if l_new_name <> 'Timmothy'  then
				l_fail = true;
				l_data = '{"status":"error","msg":"az_update_user failed to update name."}';
			end if;
		EXCEPTION
			WHEN others THEN
				l_fail = true;
				l_data = '{"status":"error","msg":"Error occured on az_update_user"}';
		END;
	END IF;

	delete 
		from "az_user"
		where "email" like '%@test-example....'
		;

	RETURN l_data;
END;
$$ LANGUAGE plpgsql;


DELETE from "az_output";

select az_test_01() ;

select "msg" from "az_output" order by "seq";
DELETE from "az_output";

