# pgcp - PostgreSQL table copy


Written in GoLang, this utility copies some or all rows of a Postgres database table as INSERT or UPDATE statements.  All columns are included in the INSERT or UPDATE statements.

A couple of binaries to save you the effort:
[Mac](https://github.com/joncrlsn/pgcp/raw/master/bin-osx/pgcp "OSX version")  [Linux](https://github.com/joncrlsn/pgcp/raw/master/bin-linux/pgcp "Linux version")

## usage

	pgcp [database flags] <genType> <tableName> [idColumn] <whereClause>


database flags | Explanation 
-------------: | -------------
  -U           | postgres user   (matches psql flag)
  -h           | database host -- default is localhost (matches psql flag)
  -p           | port.  defaults to 5432 (matches psql flag)
  -d           | database name (matches psql flag)
  -pw          | password for the postgres user<br>(if not provided then you'll be prompted)


Argument            | Explanation 
--------:           | -------------
&lt;genType&gt;     | type of SQL to generate: INSERT or UPDATE.<br/>(case insensitive)
&lt;tableName&gt;   | name of table to be outputted (fully or partially)
\[idColumn\]        | only used when genType is "update"
&lt;whereClause&gt; | specifies which rows to copy.  example:<br> "WHERE user_id < 100 AND username IS NOT NULL"

#### Database connection information can be specified in up to three ways:

  * Environment variables (keeps you from typing them in often)
  * Program flags (overrides environment variables.  See above)
  * ~/.pgpass file (may contain password for the previously specified user)
  * Note that if password is not specified, you will be prompted.

#### Optional database environment variables

Name       | Explanation
---------  | -----------
PGHOST     | host name where database is running (matches psql)
PGPORT     | port database is listening on (default is 5432) (matches psql)
PGDATABASE | name of database you want to copy (matches psql)
PGUSER     | user in postgres you'll be executing the queries as (matches psql)
PGPASSWORD | password for the user (matches psql)
PGOPTION   | one or more database options (like sslmode=disable)
