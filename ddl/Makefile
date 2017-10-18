all: sql

sql: t_az_user.sql s_login.sql s_signup.sql t_auth_token.sql s_upd_user.sql test.sql

t_az_user.sql: t_az_user.m4.sql
	/usr/local/bin/m4 -P t_az_user.m4.sql >t_az_user.sql

t_auth_token.sql: t_auth_token.m4.sql
	/usr/local/bin/m4 -P t_auth_token.m4.sql >t_auth_token.sql

s_login.sql: s_login.m4.sql
	/usr/local/bin/m4 -P s_login.m4.sql >s_login.sql

s_signup.sql: s_signup.m4.sql
	/usr/local/bin/m4 -P s_signup.m4.sql >s_signup.sql

s_upd_user.sql: s_upd_user.m4.sql
	/usr/local/bin/m4 -P s_upd_user.m4.sql >s_upd_user.sql

test.sql: test.m4.sql
	/usr/local/bin/m4 -P test.m4.sql >test.sql


