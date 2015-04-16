package main

import "fmt"

import (
	"database/sql"
	"github.com/joncrlsn/misc"
	"github.com/joncrlsn/pgutil"
	_ "github.com/lib/pq"
	flag "github.com/ogier/pflag"
	"log"
	"os"
	"strings"
	"time"
)

const isoFormat = "2006-01-02T15:04:05.000-0700"
const version = "1.0.6"

// out defaults to standard out, but can be overwritten with a traditional file
// via the -o flag
var out *os.File = os.Stdout

func print(a ...interface{}) (int, error) {
	return fmt.Fprint(out, a...)
}

func println(a ...interface{}) (int, error) {
	return fmt.Fprintln(out, a...)
}

func printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(out, format, a...)
}

/*
 * Queries the given table name and copies the column values to either an INSERT statement or
 * an UPDATE statement.
 *
 * Example: pgcp -U myuser -d mydb INSERT t_user "WHERE user_id < 4"
 */
func main() {

	var outputFileName string
	flag.StringVarP(&outputFileName, "output-file", "o", "", "Sends output to a file")

	dbInfo := pgutil.DbInfo{}
	verFlag, helpFlag := dbInfo.Populate()

	if verFlag {
		fmt.Fprintf(os.Stderr, "%s version %s\n", os.Args[0], version)
		fmt.Fprintln(os.Stderr, "Copyright (c) 2015 Jon Carlson.  All rights reserved.")
		fmt.Fprintln(os.Stderr, "Use of this source code is governed by the MIT license")
		fmt.Fprintln(os.Stderr, "that can be found here: http://opensource.org/licenses/MIT")
		os.Exit(1)
	}

	if helpFlag {
		usage()
	}

	// Remaining args:
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Missing some arguments")
		usage()
	}

	// genType
	genType := strings.ToUpper(args[0])
	if genType != "INSERT" && genType != "UPDATE" {
		fmt.Fprintf(os.Stderr, "Invalid generation type: %s.  Requires either INSERT or UPDATE\n", genType)
		usage()
	}

	if len(outputFileName) > 0 {
		var err error
		out, err = os.Create(outputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open file: %s. \n", outputFileName)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}

	// tableName
	tableName := args[1]

	// idColumn (UPDATE only) and whereClause
	whereClause := ""
	idCol := ""
	if genType == "INSERT" {
		if len(args) > 2 {
			whereClause = args[2]
		}
	} else {
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "UPDATE requires an idColumn.")
			usage()
		}
		idCol = args[2]
		if len(args) > 3 {
			whereClause = args[3]
		}
	}

	if len(whereClause) == 0 {
		// Make sure user intended there to be no where clause
		if !misc.PromptYesNo("Did you intend to have no where clause?", true) {
			os.Exit(1)
		}
	}

	db, err := dbInfo.Open()
	check("opening database", err)

	query := "SELECT * FROM " + tableName + " " + whereClause
	printf("-- Creating %s(s) from query: %s\n", genType, query)
	rowChan, columnNames := querySqlValues(db, query)

	for row := range rowChan {
		if genType == "INSERT" {
			generateInsert(tableName, row, columnNames)
		} else {
			generateUpdate(tableName, row, idCol)
		}
	}
}

func generateInsert(tableName string, row map[string]string, colNames []string) {
	printf("INSERT INTO %s (", tableName)
	first := true
	for _, name := range colNames {
		if !first {
			print(", ")
		}
		printf(name)
		first = false
	}
	print(") VALUES (")
	first = true
	for _, name := range colNames {
		if !first {
			print(", ")
		}
		v := row[name]
		printf(v)
		first = false
	}
	println(");")
}

func generateUpdate(tableName string, row map[string]string, idCol string) {
	printf("UPDATE %s SET ", tableName)
	idVal := ""
	idColFound := false
	first := true
	for k, v := range row {
		if k == idCol {
			idVal = v
			idColFound = true
		} else {
			if !first {
				print(", ")
			}
			printf("%s=%s", k, v)
			first = false
		}
	}
	if !idColFound {
		log.Fatalf("\nid column not found: %s\n", idCol)
		os.Exit(1)
	}
	printf(" WHERE %s=%s", idCol, idVal)
	println(";")
}

/*
 * Returns row maps (keyed by the column name) in a channel.
 * Dynamically converts each column value to a SQL string value.
 * See http://stackoverflow.com/questions/23507531/is-golangs-sql-package-incapable-of-ad-hoc-exploratory-queries
 */
func querySqlValues(db *sql.DB, query string) (chan map[string]string, []string) {
	rowChan := make(chan map[string]string)

	rows, err := db.Query(query)
	check("running query", err)
	columnNames, err := rows.Columns()
	check("getting column names", err)

	go func() {

		defer rows.Close()

		vals := make([]interface{}, len(columnNames))
		valPointers := make([]interface{}, len(columnNames))
		// Copy
		for i := 0; i < len(columnNames); i++ {
			valPointers[i] = &vals[i]
		}

		for rows.Next() {
			err = rows.Scan(valPointers...)
			check("scanning a row", err)

			row := make(map[string]string)
			// Convert each cell to a SQL-valid string representation
			for i, valPtr := range vals {
				//println(reflect.TypeOf(valPtr))
				switch valueType := valPtr.(type) {
				case nil:
					row[columnNames[i]] = "null"
				case []uint8:
					row[columnNames[i]] = "'" + strings.Replace(string(valPtr.([]byte)), "'", "''", -1) + "'"
				case string:
					row[columnNames[i]] = "'" + strings.Replace(valPtr.(string), "'", "''", -1) + "'"
				case int64:
					row[columnNames[i]] = fmt.Sprintf("%d", valPtr)
				case float64:
					row[columnNames[i]] = fmt.Sprintf("%f", valPtr)
				case bool:
					row[columnNames[i]] = fmt.Sprintf("%t", valPtr)
				case time.Time:
					row[columnNames[i]] = "'" + valPtr.(time.Time).Format(isoFormat) + "'"
				case fmt.Stringer:
					row[columnNames[i]] = fmt.Sprintf("%v", valPtr)
				default:
					row[columnNames[i]] = fmt.Sprintf("%v", valPtr)
					println("-- Warning, column %s is an unhandled type: %v", columnNames[i], valueType)
				}
			}
			rowChan <- row
		}
		close(rowChan)
	}()
	return rowChan, columnNames
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [database flags] <genType> <tableName> <whereClause> \n", os.Args[0])
	fmt.Fprintln(os.Stderr, `
Copies table data as either INSERT or UPDATE statements.

Program flags are:
  -V, --version      : prints the version of pgcp being run
  -?, --help         : prints a summary of the commands accepted by pgcp
  -U, --user         : user in postgres to execute the commands
  -h, --host         : host name where database is running (default is localhost)
  -p, --port         : port database is listening on (default is 5432)
  -d, --dbname       : database name
  -O, --options      : postgresql connection options (like sslmode=disable)
  -w, --no-password  : Never issue a password prompt
  -W, --password     : Force a password prompt
  -o, --output-file  : Sends output to the given file

<genType>     : type of SQL to generate: insert, update

<tableName>   : name of table to be copied (fully or partially)

<whereClause> : specifies which rows to copy. Example: WHERE user_id < 100 AND username IS NOT NULL

Database connection information can be specified in two ways:
  * Environment variables
  * Program flags (overrides environment variables.  See above)
  * ~/.pgpass file (for the password)
  * Note that if password is not specified, you will be prompted.

Optional Environment variables (if program flags are not desirable):
  PGHOST     : host name where database is running (default is localhost)
  PGPORT     : port database is listening on (default is 5432)
  PGDATABASE : name of database you want to copy
  PGUSER     : user in postgres you'll be executing the queries as
  PGPASSWORD : password for the postgres user
  PGOPTION   : options (like sslmode=disable)
`)

	os.Exit(2)
}

func check(msg string, err error) {
	if err != nil {
		log.Fatal("Error "+msg, err)
	}
}
