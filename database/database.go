package database

// This package encapsulates a lot of the boilerplate associated with the database/sql package.
// Errors are handled internally unless they indicate a failure.

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	// This blank include forces Init() and does not need to be in main.
	_ "github.com/knousere/web-service-commons/go-sql-driver/mysql"
	"github.com/knousere/web-service-commons/utils"
)

// DBConnection is a container for database connection parameters.
type DBConnection struct {
	db           *sql.DB
	Engine       string // mysql
	Scheme       string // https
	Host         string // 123.123.123.123:3306
	Schema       string // database schema
	User         string // database user name
	PasswordPath string // relative path of password file
}

// AppDb application database instance
var AppDb *DBConnection

// LogDb logging database instance
var LogDb *DBConnection

// Open up a connection to the database and ping it to make sure that it actually works
func (dbConn *DBConnection) Open() error {
	strPassword, err := dbConn.readDbPassword()
	if err != nil {
		fmt.Println("database.Open failed on readDbPassword", err.Error())
		return err
	}
	strConnection := dbConn.GetConnectionString(strPassword)
	// utils.Info.Println(strConnection)
	dbConn.db, err = sql.Open(dbConn.Engine, strConnection)
	if err != nil {
		fmt.Println("database.Open failed on sql.Open", err.Error())
		return err
	}

	err = dbConn.db.Ping()
	if err != nil {
		fmt.Println("database.Open failed on db.Ping", err.Error())
		return err
	}

	return nil
}

// Close database connection explicitly.
func (dbConn *DBConnection) Close() {
	if dbConn.db != nil {
		dbConn.db.Close()
	}
}

func (dbConn *DBConnection) readDbPassword() (string, error) {
	passBytes, err := ioutil.ReadFile(dbConn.PasswordPath)
	var strPassword string
	if err == nil {
		strPassword = string(passBytes)
		strPassword = utils.CleanPassword(strPassword)
	}
	return strPassword, err
}

// GetConnectionString builds a database connection string.
func (dbConn *DBConnection) GetConnectionString(strPassword string) string {
	u := url.URL{
		Scheme: dbConn.Scheme,
		Host:   fmt.Sprintf("tcp(%s)", dbConn.Host),
		Path:   dbConn.Schema,
		User:   url.UserPassword(dbConn.User, strPassword),
	}
	query := u.Query()
	query.Set("parseTime", "true")
	query.Set("autocommit", "true")
	query.Set("collation", "utf8mb4_unicode_ci")
	u.RawQuery = query.Encode()
	strConnection := u.String()
	if strings.HasPrefix(strConnection, "//") {
		return strConnection[2:]
	}
	return strConnection
}
