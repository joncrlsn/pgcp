package main

import "fmt"

//import "reflect"
import "time"
import "log"
import flag "github.com/ogier/pflag"
import "os"
import "strings"
import "database/sql"
import _ "github.com/lib/pq"
import "github.com/joncrlsn/pgutil"
import "github.com/joncrlsn/misc"

const isoFormat = "2006-01-02T15:04:05.000-0700"
const version = "1.0.4"

/*
 * Queries the given table name and copies the column values to either an INSERT statement or
 * an UPDATE statement.
 *
 * Example: pgcp -U myuser -d mydb INSERT t_user "WHERE user_id < 4"
 */
func main() {

	var verFlag bool
	var helpFlag bool
	flag.BoolVarP(&verFlag, "version", "V", false, "Displays version information")
	flag.BoolVarP(&helpFlag, "help", "?", false, "Displays usage help")

	dbInfo := pgutil.DbInfo{}
	dbInfo.Populate()

	if verFlag {
		fmt.Printf("%s - version %s\n", os.Args[0], version)
		os.Exit(0)
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
	fmt.Printf("-- Creating %s(s) from query: %s\n", genType, query)
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
	fmt.Printf("INSERT INTO %s (", tableName)
	first := true
	for _, name := range colNames {
		if !first {
			fmt.Print(", ")
		}
		fmt.Printf(name)
		first = false
	}
	fmt.Print(") VALUES (")
	first = true
	for _, name := range colNames {
		if !first {
			fmt.Print(", ")
		}
		v := row[name]
		fmt.Printf(v)
		first = false
	}
	fmt.Println(");")
}

func generateUpdate(tableName string, row map[string]string, idCol string) {
	fmt.Printf("UPDATE %s SET ", tableName)
	idVal := ""
	idColFound := false
	first := true
	for k, v := range row {
		if k == idCol {
			idVal = v
			idColFound = true
		} else {
			if !first {
				fmt.Print(", ")
			}
			fmt.Printf("%s=%s", k, v)
			first = false
		}
	}
	if !idColFound {
		log.Fatalf("\nid column not found: %s\n", idCol)
		os.Exit(1)
	}
	fmt.Printf(" WHERE %s=%s", idCol, idVal)
	fmt.Println(";")
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
				//fmt.Println(reflect.TypeOf(valPtr))
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
					fmt.Println("-- Warning, column %s is an unhandled type: %v", columnNames[i], valueType)
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
