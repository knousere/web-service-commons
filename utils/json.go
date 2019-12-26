package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// UnmarshalRequest unmarshals an http request body as json into
// a pointer to struct as interface and closes the request body.
// Any error should be reported as the responsibility of the requester.
func UnmarshalRequest(r *http.Request, retJSON interface{}) error {
	defer r.Body.Close() // allow the connection to be reused

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Warning.Println("failed on ioutil.ReadAll", err.Error())
		return err
	}

	err = json.Unmarshal(body, retJSON)
	if err != nil {
		Warning.Println("failed on json.Unmarshal", err.Error())
		return err
	}

	return nil
}

// UnmarshalResponse unmarshals an http response body as json into
// a pointer to struct as interface and closes the response body.
// Any error should be reported as the responsibility of the requester.
func UnmarshalResponse(r *http.Response, retJSON interface{}) error {
	defer r.Body.Close() // allow the connection to be reused

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Warning.Println("failed on ioutil.ReadAll", err.Error())
		return err
	}

	err = json.Unmarshal(body, retJSON)
	if err != nil {
		Warning.Println("failed on json.Unmarshal", err.Error())
		return err
	}

	return nil
}

// JSONStringify marshals an object into a json string.
// It intentionally does not pass an error back.
func JSONStringify(v interface{}) string {
	j, err := json.Marshal(v)
	if err != nil {
		Warning.Println("JSONStringify error", err.Error())
		return ""
	}

	str := string(j)
	return str
}

// JSONPrettyPrint formats a json string into an indented form for display.
func JSONPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

// SQ wraps target string with single quotes.
func SQ(strTarget string) string {
	strSQ := "'"
	return strSQ + strTarget + strSQ
}

// DQ wraps target string with double quotes.
func DQ(strTarget string) string {
	strDQ := "\""
	return strDQ + strTarget + strDQ
}

// BQ wraps target string with back quotes.
func BQ(strTarget string) string {
	strBQ := "`"
	return strBQ + strTarget + strBQ
}

// CB wraps target string with curly braces.
func CB(strTarget string) string {
	strCB1 := "{"
	strCB2 := "}"
	return strCB1 + strTarget + strCB2
}

// SB wraps target string with square braces.
func SB(strTarget string) string {
	strSB1 := "["
	strSB2 := "]"
	return strSB1 + strTarget + strSB2
}

// JSONKeyString formats a json key value pair as string {"key": "value"}.
func JSONKeyString(strKey string, strValue string) string {
	return fmt.Sprintf(`{"%s": "%s"}`, strKey, strValue)
}

// JSONKeyInt formats a json key value pair as int {"key": intValue}.
func JSONKeyInt(strKey string, intValue int) string {
	return fmt.Sprintf(`{"%s": %d}`, strKey, intValue)
}

// JSONKeyBool formats json key value pair as boolean {"key": <true/false>}.
func JSONKeyBool(strKey string, bValue bool) string {
	return fmt.Sprintf(`{"%s": %t}`, strKey, bValue)
}

// JSONKeyArrayToString converts a key value array to a segment of a json dictionary.
func JSONKeyArrayToString(keyVals []KeyValuePair) string {
	strPairs := make([]string, 0, 30)
	for _, keyVal := range keyVals {
		strPair := fmt.Sprintf(`"%s": "%s"`, keyVal.Key, keyVal.Value)
		strPairs = append(strPairs, strPair)
	}
	return strings.Join(strPairs, ",")
}

// JSONMapToString converts a key value map to a segment of a json dictionary.
func JSONMapToString(keyVals map[string]string) string {
	strPairs := make([]string, 0, 30)
	for key, val := range keyVals {
		strPair := fmt.Sprintf(`"%s": "%s"`, key, val)
		strPairs = append(strPairs, strPair)
	}
	return strings.Join(strPairs, ",")
}
