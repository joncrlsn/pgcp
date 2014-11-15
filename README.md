# pgcp - PostgreSQL table copy


Written in GoLang, this utility copies some or all rows of a Postgres database table as INSERT or UPDATE statements.  All columns are included in the INSERT or UPDATE statements.

Suggestions to make this more useful and "idiomatic Go" will be appreciated.

### download
[osx64](https://github.com/joncrlsn/pgcp/raw/master/bin-osx64/pgcp "OSX 64-bit version") 
[osx32](https://github.com/joncrlsn/pgcp/raw/master/bin-osx32/pgcp "OSX version")
[linux64](https://github.com/joncrlsn/pgcp/raw/master/bin-linux64/pgcp "Linux 64-bit version")
[linux32](https://github.com/joncrlsn/pgcp/raw/master/bin-linux32/pgcp "Linux version")
[win64](https://github.com/joncrlsn/pgcp/raw/master/bin-win64/pgcp.exe "Windows 64-bit version")
[win32](https://github.com/joncrlsn/pgcp/raw/master/bin-win32/pgcp.exe "Windows version")

### usage
	pgcp [database flags] <genType> <tableName> [idColumn] <whereClause>

### examples
	pgcp -U dbuser -h 10.10.41.55 -d userdb INSERT users         "where user_id > 10"
	pgcp -U dbuser -h 10.10.41.55 -d userdb UPDATE users user_id "where user_id > 10"

### flags
database flag | Explanation 
------------: | -------------
  -U          | postgres user   (matches psql flag)
  -h          | database host -- default is localhost (matches psql flag)
  -p          | port.  defaults to 5432 (matches psql flag)
  -d          | database name (matches psql flag)
  -pw         | password for the postgres user<br>(if not provided then you'll be prompted)


Argument            | Explanation 
--------:           | -------------
&lt;genType&gt;     | type of SQL to generate: INSERT or UPDATE.<br/>(case insensitive)
&lt;tableName&gt;   | name of table to be outputted (fully or partially)
\[idColumn\]        | only used when genType is UPDATE
&lt;whereClause&gt; | specifies which rows to copy.  example:<br> "WHERE user_id < 100 AND username IS NOT NULL"

### database connection information can be specified in up to three ways:

  * Environment variables (keeps you from typing them in often)
  * Program flags (overrides environment variables.  See above)
  * ~/.pgpass file (may contain password for the previously specified user)
  * Note that if password is not specified, you will be prompted.

### optional database environment variables

Name       | Explanation
---------  | -----------
PGHOST     | host name where database is running (matches psql)
PGPORT     | port database is listening on (default is 5432) (matches psql)
PGDATABASE | name of database you want to copy (matches psql)
PGUSER     | user in postgres you'll be executing the queries as (matches psql)
PGPASSWORD | password for the user (matches psql)
PGOPTION   | one or more database options (like sslmode=disable)
