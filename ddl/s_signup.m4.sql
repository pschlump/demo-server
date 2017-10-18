
DROP FUNCTION az_signup ( p_email varchar, p_password varchar, p_first_name varchar, p_last_name varchar);

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
			l_auth_token = uuid_generate_v4();
			insert into "az_auth_token" (
				  "auth_token"	
				, "user_id"	
			) values (
				  l_auth_token
				, l_id
			);
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

-- delete from "az_user";
-- delete from "az_auth_token";
-- select az_signup( 'a@b.c', 'abc', 'a', 'b');
-- select * from "az_user";
-- select * from "az_auth_token";

