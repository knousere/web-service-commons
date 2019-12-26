package utils

import (
	"bufio"
	"os"
	"strings"
)

// Params is exported as a map of command line parameters.
// By convention the keys are lower case.
var Params map[string]string

// ReadParams reads an initialization parameter file into the Params map.
func ReadParams(strParamPath string) error {
	Params = make(map[string]string)
	f, err := os.Open(strParamPath)
	if err != nil {
		Warning.Println("failed to open", strParamPath, err.Error())
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		Trace.Println(scanner.Text())
		setParam(scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		Warning.Println("reading param file:", strParamPath, err.Error())
		return err
	}
	Trace.Printf("Params len %d\n", len(Params))
	return nil
}

// setParam parses a line from a parameter file into the Params map
func setParam(strText string) {
	strText = strings.TrimSpace(strText)
	if strText != "" && !strings.HasPrefix(strText, "#") {
		s := strings.SplitN(strText, ":", 2)
		if len(s) == 2 {
			strKey := strings.ToLower(strings.TrimSpace(s[0]))
			strValue := strings.TrimSpace(s[1])
			Params[strKey] = strValue
		}
	}
}
