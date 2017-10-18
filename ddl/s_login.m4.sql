
drop FUNCTION az_login(p_email varchar, p_password varchar);

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

	IF not l_fail THEN
		select  "az_user"."id"
		into  l_id
			from "az_user" as "az_user" left join "az_auth_token" as "az_auth_token" on "az_auth_token"."user_id" = "az_user"."id"
			where "az_user"."email" = p_email
			  and "az_user"."password" = crypt(p_password, "az_user"."password")
			;

		IF not found THEN
			l_data = '{ "status":"error", "code":"401", "msg":"Invalid username or password." }';
			l_fail = true;
		END IF;
	END IF;

	IF not l_fail THEN
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
	END IF;

	RETURN l_data;
END;
$$ LANGUAGE plpgsql;

-- select az_login('a@b.c', 'abc') ;
-- select * from "az_user";
-- select * from "az_auth_token";

