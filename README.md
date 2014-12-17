# pgcp - PostgreSQL table copy

Written in GoLang, this utility copies some or all rows of a Postgres database table as INSERT or UPDATE statements.  All columns are included in the INSERT or UPDATE statements.

Suggestions and modifications to make this more useful and "idiomatic Go" will be appreciated.

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
	pgcp -U dbuser -h 10.10.41.55 -d userdb INSERT users         "where user_id < 10"
	pgcp -U dbuser -h 10.10.41.55 -d userdb UPDATE users user_id "where user_id < 10"

#### options/flags (these mostly match psql arguments):
program flag/option  | explanation
-------------------: | -------------
  -V, --version      | prints the version of pgcp being run
  -?, --help         | prints a summary of the commands accepted by pgcp
  -U, --user         | user in postgres to execute the commands
  -h, --host         | host name where database is running (default is localhost)
  -p, --port         | port database is listening on (default is 5432)
  -d, --dbname       | database name
  -O, --options      | postgresql connection options (like sslmode=disable)
  -w, --no-password  | Never issue a password prompt
  -W, --password     | Force a password prompt

argument            | explanation 
--------:           | -------------
&lt;genType&gt;     | type of SQL to generate: INSERT or UPDATE.<br/>(case insensitive)
&lt;tableName&gt;   | name of table to be outputted (fully or partially)
\[idColumn\]        | only specify when genType is UPDATE
&lt;whereClause&gt; | specifies which rows to copy.  example:<br> "WHERE user_id < 100 AND username IS NOT NULL"

### database connection options

  * Use environment variables (see table below)
  * Program flags (overrides environment variables)
  * ~/.pgpass file
  * Note that if password is not specified, you will be prompted.

### optional database environment variables

name       | explanation
---------  | -----------
PGHOST     | host name where database is running (matches psql)
PGPORT     | port database is listening on (default is 5432) (matches psql)
PGDATABASE | name of database you want to copy (matches psql)
PGUSER     | user in postgres you'll be executing the queries as (matches psql)
PGPASSWORD | password for the user (matches psql)
PGOPTION   | one or more database options (like sslmode=disable)

### todo
1. Fix bug where Ctrl-C in the password entry field messes up the console.
2. Convert positional arguments to program options/flags?  How important is this to people?
3. Add database options that may be requested by others (that fit with the purpose of this tool).