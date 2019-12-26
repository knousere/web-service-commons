package database

// Many of these wrapper functions are designed to call mysql stored procedures that return
// result codes rather than raw errors. A stored procedure cannot be called using the exec function.
// All stored procedures must return either a dataset or a result code. The clientMultiResults
// flag is turned on internally. However, multiple return sets are not supported and will be
// silently discarded if found.

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/knousere/web-service-commons/utils"
)

// Exec is a wrapper for the sql.DB.Exec function.
// WARNING: do no use this for a stored procedure.
func (dbConn *DBConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	strArgs := doTrace(query, args...)
	result, err := dbConn.db.Exec(query, args...)
	if err != nil {
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	}
	if utils.GetTrace() == 1 {
		lastInsertID, _ := result.LastInsertId()
		affectedRows, _ := result.RowsAffected()
		utils.Trace.Printf("AffectedRows=%d, LastInsertId=%d", affectedRows, lastInsertID)
	}
	return result, err
}

// argString converts a arg list into a string.
func argString(query string, args ...interface{}) string {

	argArray := make([]string, 0, 20)
	argArray = append(argArray, query)

	for _, v := range args {
		strToken := fmt.Sprintf("%v", v)
		argArray = append(argArray, strToken)
	}
	strArray := strings.Join(argArray, ", ")
	return strArray
}

// doTrace puts calling arguments into trace if Trace is set.
// return argString, empty of trace is not set
func doTrace(query string, args ...interface{}) string {
	var strArgs string
	if utils.GetTrace() == 1 {
		strArgs = argString(query, args...)
		utils.Trace.Println(strArgs)
	}
	return strArgs
}

// refreshTrace guarantees that strArgs has been initialized once.
// This prevents unecessary repeated calls to argString.
func refreshTrace(strArgs, query string, args ...interface{}) string {
	if strArgs == "" {
		strArgs = argString(query, args...)
	}
	return strArgs
}

// LogError is a wrapper for whatever logging facility is employed locally.
// The implementation line should be modified as necessary.
func (dbConn *DBConnection) LogError(err error, query string, args ...interface{}) {
	utils.Warning.Println(err.Error(), argString(query, args...))
}

// GetRows simply gets rows. Empty is OK.
// Errors are logged here so calling code need not be cluttered.
// Most impartantly this never returns ErrNoRows.
func (dbConn *DBConnection) GetRows(query string, args ...interface{}) (*sql.Rows, error) {
	strArgs := doTrace(query, args...)

	rows, err := dbConn.db.Query(query, args...)
	switch {
	case err == sql.ErrNoRows:
		return rows, nil
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
		return rows, err
	default:
		return rows, nil
	}
}

// GetOneRow is functionally identical to QueryRow().
// This isolates the get single row case for debugging purposes
func (dbConn *DBConnection) GetOneRow(query string, args ...interface{}) *sql.Row {
	_ = doTrace(query, args...)
	return dbConn.db.QueryRow(query, args...)
}

// GetPositiveInt gets one row that consists of only a positive integer.
// This is usually either a record id or a count.
// A negative value indicates a problem.
func (dbConn *DBConnection) GetPositiveInt(query string, args ...interface{}) (int, error) {
	strArgs := doTrace(query, args...)
	var intValue int
	err := dbConn.db.QueryRow(query, args...).Scan(&intValue)
	switch {
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case intValue < 0:
		utils.Warning.Printf("query returned value:%d %s", intValue, refreshTrace(strArgs, query, args...))
	}
	return intValue, err
}

// GetPositiveIntDefault  gets one row that consists of only a positive integer.
// It returns default value rather than ErrNoRows.
func (dbConn *DBConnection) GetPositiveIntDefault(intDefault int, query string, args ...interface{}) (int, error) {
	strArgs := doTrace(query, args...)
	var intValue int
	err := dbConn.db.QueryRow(query, args...).Scan(&intValue)
	switch {
	case err == sql.ErrNoRows:
		return intDefault, nil
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case intValue < 0:
		utils.Warning.Printf("query returned value:%d %s", intValue, refreshTrace(strArgs, query, args...))
	}
	return intValue, err
}

// GetRecordID returns record Id or 0 rather than ErrNoRows.
func (dbConn *DBConnection) GetRecordID(query string, args ...interface{}) (int, error) {
	strArgs := doTrace(query, args...)
	var intID int
	err := dbConn.db.QueryRow(query, args...).Scan(&intID)
	switch {
	case err == sql.ErrNoRows:
		utils.Trace.Println("record not found", refreshTrace(strArgs, query, args...))
		return 0, nil
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case intID == 0:
		utils.Trace.Println("record not found", refreshTrace(strArgs, query, args...))
	case intID < 0:
		utils.Warning.Printf("query returned id:%d %s", intID, refreshTrace(strArgs, query, args...))
	}
	return intID, err
}

// GetRecordCount returns record count. ErrNoRows is an error.
func (dbConn *DBConnection) GetRecordCount(query string, args ...interface{}) (int, error) {
	strArgs := doTrace(query, args...)
	var intCount int
	err := dbConn.db.QueryRow(query, args...).Scan(&intCount)
	switch {
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case intCount < 0:
		utils.Warning.Printf("query returned count:%d %s", intCount, refreshTrace(strArgs, query, args...))
	}
	return intCount, err
}

// GetOneString get one row that consists of only a string value.
// This is usually an external key.
func (dbConn *DBConnection) GetOneString(query string, args ...interface{}) (string, error) {
	strArgs := doTrace(query, args...)
	var strValue string
	err := dbConn.db.QueryRow(query, args...).Scan(&strValue)

	if err != nil {
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	}

	return strValue, err
}

// InsertRow inserts a row and returns the id of the new record.
// An ID should always be returned so all errors including "no rows" are an
// indication of a server problem.
func (dbConn *DBConnection) InsertRow(query string, args ...interface{}) (int, error) {
	strArgs := doTrace(query, args...)
	var id int
	err := dbConn.db.QueryRow(query, args...).Scan(&id)
	switch {
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case id < 1:
		utils.Warning.Printf("query returned id:%d %s", id, refreshTrace(strArgs, query, args...))
	}
	return id, err
}

// InsertRowResult inserts a row and returns the result code and the id of the new record.
// Result -2 typically indicates a permissions error.
// An ID should always be returned so all errors including "no rows" indicate a server problem.
func (dbConn *DBConnection) InsertRowResult(query string, args ...interface{}) (int, int, error) {
	strArgs := doTrace(query, args...)
	var result, id int
	err := dbConn.db.QueryRow(query, args...).Scan(&result, &id)
	switch {
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case id < 1:
		utils.Warning.Printf("query returned id:%d %s", id, refreshTrace(strArgs, query, args...))
	}
	return result, id, err
}

// UpdateRows updates row(s) and returns affected count.
// An affected count should always be returned so all errors including "no rows" indicate
// a server problem.
func (dbConn *DBConnection) UpdateRows(query string, args ...interface{}) (int, error) {
	strArgs := doTrace(query, args...)
	var affectedCount int
	err := dbConn.db.QueryRow(query, args...).Scan(&affectedCount)
	switch {
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case affectedCount < 0:
		utils.Warning.Printf("query returned count:%d %s", affectedCount, refreshTrace(strArgs, query, args...))
	}
	return affectedCount, err
}

// UpdateRowsWithDeadlock updates row(s) and returns count of affected rows and deadlock flag.
// The query is expected to be a stored procedure.
func (dbConn *DBConnection) UpdateRowsWithDeadlock(query string, args ...interface{}) (int, bool) {
	strArgs := doTrace(query, args...)
	var affectedCount int
	var bDeadlock bool
	err := dbConn.db.QueryRow(query, args...).Scan(&affectedCount, &bDeadlock)
	if err != nil {
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	}
	return affectedCount, bDeadlock
}

// UpdateRowsResult updates row(s) and returns result code and affected count if there was no error in SP.
// A result code should always be returned so all errors including "no rows" indicate
// a server problem.
func (dbConn *DBConnection) UpdateRowsResult(query string, args ...interface{}) (int, int, error) {
	strArgs := doTrace(query, args...)
	var result, affectedCount int
	err := dbConn.db.QueryRow(query, args...).Scan(&result, &affectedCount)
	switch {
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
		result = -1
	case result < 0:
		utils.Warning.Printf("query returned result:%d %s", result, refreshTrace(strArgs, query, args...))
	case affectedCount < 0:
		utils.Warning.Printf("query returned count:%d %s", affectedCount, refreshTrace(strArgs, query, args...))
	}
	return result, affectedCount, err
}

// DeleteRows deletes row(s) and returns affected count.
// An affected count should always be returned so all errors including "no rows" indicate
// a server problem.
// Yes, this looks identical to UpdateRows. It is segregated to make debugging more clear.
func (dbConn *DBConnection) DeleteRows(query string, args ...interface{}) (int, error) {
	strArgs := doTrace(query, args...)
	var affectedCount int
	err := dbConn.db.QueryRow(query, args...).Scan(&affectedCount)
	switch {
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case affectedCount < 0:
		utils.Warning.Printf("query returned count:%d %s", affectedCount, refreshTrace(strArgs, query, args...))
	}
	return affectedCount, err
}

// DeleteRowsResult deletes row(s) and returns result code as well as an affected count.
// The query is expected to be a stored procedure.
// An affected count should always be returned so all errors including "no rows" indicate
// a server problem.
func (dbConn *DBConnection) DeleteRowsResult(query string, args ...interface{}) (int, int, error) {
	strArgs := doTrace(query, args...)
	var result, affectedCount int
	err := dbConn.db.QueryRow(query, args...).Scan(&result, &affectedCount)
	switch {
	case err != nil:
		utils.Warning.Println(err.Error(), refreshTrace(strArgs, query, args...))
	case result < 0:
		utils.Warning.Printf("query returned result:%d %s", result, refreshTrace(strArgs, query, args...))
	case affectedCount < 0:
		utils.Warning.Printf("query returned count:%d %s", affectedCount, refreshTrace(strArgs, query, args...))
	}
	return result, affectedCount, err
}
