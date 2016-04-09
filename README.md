# pgcp - PostgreSQL table copy

Written in GoLang, this utility copies some or all rows of a Postgres database table as INSERT or UPDATE statements.  All columns are included in the INSERT or UPDATE statements.

Suggestions and modifications to make this more useful and "idiomatic Go" will be appreciated.

### download
[osx](https://github.com/joncrlsn/pgcp/raw/master/bin-osx/pgcp "OSX version") 
[linux](https://github.com/joncrlsn/pgcp/raw/master/bin-linux/pgcp "Linux version")
[windows](https://github.com/joncrlsn/pgcp/raw/master/bin-win/pgcp.exe "Windows version")

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
  -w, --no-password  | Never issue a db password prompt
  -W, --password     | Force a db password prompt
  -o, --output-file  | Send output to the given file  

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
1. ~~Fix bug where Ctrl-C in the password entry field messes up the console.~~ fixed in version 1.0.5
1. ~~Fix -? and -V flags that are not working.~~ Fixed in version 1.0.6
1. ~~Add --output-file (-o) flag~~ Added in version 1.0.6
2. Convert positional arguments to program options/flags?  How important is this to people?
3. Improve the accuracy of parsing ~/.pgpass
