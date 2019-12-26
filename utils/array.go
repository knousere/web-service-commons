package utils

import (
	"strconv"
	"strings"
)

// KeyValuePair is a container for a key value pair for JSON manipulation.
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// FindKeyValue returns true and the value if the KeyValuePair array contains the key.
func FindKeyValue(array []KeyValuePair, key string) (bool, string) {
	for _, a := range array {
		if a.Key == key {
			return true, a.Value
		}
	}
	return false, ""
}

// IntArray2String converts an array of int to a comma delimited string up to max length.
// Return remainder int array if necessary.
func IntArray2String(array []int, intMaxLen int) (string, []int) {
	var strList string
	var intLength int
	strInts := make([]string, 0, 30)
	var i int

	for i = 0; i < len(array); i++ {
		intItem := array[i]
		strItem := strconv.Itoa(intItem)
		if intMaxLen > -1 && intLength+len(strItem) > intMaxLen {
			break
		} else {
			strInts = append(strInts, strItem)
			intLength += len(strItem) + 1
		}
	}
	//Trace.Printf("intLength=%d, i=%d\n", intLength, i)
	strList = strings.Join(strInts, ",")
	intListRet := array[i:]
	return strList, intListRet
}

// String2IntArray converts a delimited string to an array of int.
// Zeros, errors and non-digits are ignored.
func String2IntArray(strRaw string, strDelim string) []int {
	strArray := strings.Split(strRaw, strDelim)

	intArray := make([]int, 0, len(strArray))

	for _, strToken := range strArray {
		strToken = strings.TrimSpace(strToken)
		intToken, _ := strconv.Atoi(strToken)
		if intToken != 0 {
			intArray = append(intArray, intToken)
		}
	}
	return intArray
}

// ContainsString returns true and the element if the string array contains the item.
// Match case insensitive.
func ContainsString(array []string, item string) (bool, string) {
	itemLC := strings.ToLower(item)
	for _, element := range array {
		elementLC := strings.ToLower(element)
		if elementLC == itemLC {
			return true, element
		}
	}
	return false, ""
}

// InsertHeadInt inserts int as head of the array.
func InsertHeadInt(array []int, item int) []int {
	retList := make([]int, 0, len(array)+1)
	switch len(array) {
	case 0:
		retList = append(retList, item)
	case 1:
		retList = append(retList, item)
		if array[0] != item {
			retList = append(retList, array[0])
		}
	default: // len > 1
		if array[0] == item {
			retList = array
		} else {
			retList = append(retList, item)
			retList = append(retList, array...)
		}
	}
	return retList
}

// RemoveDupeInt removes second occurance of a dupe from sorted int array.
func RemoveDupeInt(array []int) ([]int, bool) {
	bFound := false
	retList := array[:]
	switch len(array) {
	case 0, 1:
		// no dupes. done
	default: // len > 1
		for i, trial := range array {
			if i+1 == len(array) {
				break // no dupes, done
			}
			a := array[i+1:]
			j := FindInt(a, trial)
			if j >= 0 {
				retList = append(retList[:i+j+1], retList[i+j+2:]...)
				bFound = true
				break // dupe found and parsed out
			}
		}
	}
	return retList, bFound
}

// FindInt returns index if the integer array contains the item else -1.
func FindInt(array []int, item int) int {
	for i, a := range array {
		if a == item {
			return i
		}
	}
	return -1
}

// ContainsInt returns true if the integer array contains the item.
func ContainsInt(array []int, item int) bool {
	for _, a := range array {
		if a == item {
			return true
		}
	}
	return false
}

// ContainsDupeInt returns an element if it appears more than once
// in a sorted int array else 0.
func ContainsDupeInt(array []int) int {
	if len(array) > 1 {
		for i, trial := range array {
			if i+1 == len(array) {
				return 0
			}
			a := array[i+1:]
			if ContainsInt(a, trial) == true {
				return trial
			}
		}
	}
	return 0
}
