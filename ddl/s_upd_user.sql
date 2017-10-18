

DROP FUNCTION az_update_user ( p_auth_token varchar,  p_first_name varchar, p_last_name varchar);

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

-- select az_update_user( '847e3746-52ee-4feb-b59a-e5601cf4903e', 'aaaa', 'bbbb');
-- select * from "az_user";
-- select * from "az_auth_token";

