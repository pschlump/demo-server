
--
-- Remember to (on linux) - use 9.5, 9.6 or 10.0 depending on version of database.
--
-- 		$ sudo apt-get install postgresql-contrib-9.5
--
-- Before running this.
--
-- Must run as "postgres" user
--       ALTER ROLE pschlump SUPERUSER;
--

CREATE EXTENSION "uuid-ossp";
CREATE EXTENSION pgcrypto;

