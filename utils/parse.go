package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const timeFormat = "2006-01-02 15:04:05"

// ParseTimestamp parses a string into a timestamp string in standard format.
// This is usually used to prove 'since' parameters.
// Return true and the formatted string if successful.
func ParseTimestamp(strTrial string) (bool, string) {

	bSuccess := false
	var strOut string
	var err error
	var tm time.Time
	var unixSecs int64

	//Trace.Println("strTrial", strTrial)
	strTrial = strings.TrimSpace(strTrial)

	switch {
	case strTrial == "":
		bSuccess = false
	case strings.Contains(strTrial, "-"):
		tm, err = time.Parse(timeFormat, strTrial)
		if err == nil {
			bSuccess = true
		}
	case strings.Contains(strTrial, "."):
		strSegs := strings.Split(strTrial, ".")
		unixSecs, err = strconv.ParseInt(strSegs[0], 10, 64)
		if err == nil {
			tm = time.Unix(unixSecs, 0)
			bSuccess = true
		}
	default:
		unixSecs, err = strconv.ParseInt(strTrial, 10, 64)
		if err == nil {
			tm = time.Unix(unixSecs, 0)
			bSuccess = true
		}
	}
	//Trace.Printf("bSuccess %v\n", bSuccess)
	if bSuccess == true {
		strOut = tm.Format(timeFormat)
	}
	//Trace.Println("strOut", strOut)
	return bSuccess, strOut
}

// ParseHashTags parses a string array of hashtags from strMessage which is interpreted as utf8 runes.
// A hashtag begins with # and ends with whitespace, punctuation, end of string, @ or #.
func ParseHashTags(strMessage string) []string {
	bInTag := false
	runes := make([]string, 0, 30)
	strTag := ""
	buf := make([]byte, 4)
	strTags := make([]string, 0, 30)
	var n int

	for _, c := range strMessage {
		switch {
		case c == '#':
			if bInTag == true {
				strTag = strings.Join(runes, "")
				strTags = append(strTags, strTag)
			}
			runes = make([]string, 0, 30)
			bInTag = true
		case unicode.IsSpace(c) || unicode.IsPunct(c) || c == '@':
			if bInTag == true {
				strTag = strings.Join(runes, "")
				strTags = append(strTags, strTag)
				runes = make([]string, 0, 30)
				bInTag = false
			}
		default:
			if bInTag == true {
				n = utf8.EncodeRune(buf, c)
				runes = append(runes, string(buf[:n]))
			}
		}
	}
	if bInTag == true {
		strTag = strings.Join(runes, "")
		strTags = append(strTags, strTag)
	}

	return strTags
}

// ParseMentions parses a string array of mentions from strMessage which is interpreted as utf8 runes.
// A mention begins with @ and ends with whitespace, punctuation, end of string, @ or #.
func ParseMentions(strMessage string) []string {
	bInTag := false
	runes := make([]string, 0, 30)
	strTag := ""
	buf := make([]byte, 4)
	strTags := make([]string, 0, 30)
	var n int

	for _, c := range strMessage {
		switch {
		case c == '@':
			if bInTag == true {
				strTag = strings.Join(runes, "")
				strTags = append(strTags, strTag)
			}
			runes = make([]string, 0, 30)
			bInTag = true
		case unicode.IsSpace(c) || unicode.IsPunct(c) || c == '#':
			if bInTag == true {
				strTag = strings.Join(runes, "")
				strTags = append(strTags, strTag)
				runes = make([]string, 0, 30)
				bInTag = false
			}
		default:
			if bInTag == true {
				n = utf8.EncodeRune(buf, c)
				runes = append(runes, string(buf[:n]))
			}
		}
	}
	if bInTag == true {
		strTag = strings.Join(runes, "")
		strTags = append(strTags, strTag)
	}

	return strTags
}

// PreviewString parses off the first 5 words of the message followed by ...
func PreviewString(strMessage string) string {
	parts := strings.Split(strMessage, " ")
	strPreview := strMessage
	if len(parts) > 5 {
		strPreview = strings.Join(parts[:5], " ") + "..."
	}
	return strPreview
}

const emailPattern = `(\w[+-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`

// ValidEmail returns true if the email matches the standard pattern.
func ValidEmail(strEmail string) bool {
	result, _ := regexp.Compile(emailPattern)
	return result.MatchString(strEmail)
}

// CleanFacebook returns Facebook name as lower case trimmed string.
func CleanFacebook(strFacebook string) string {
	return strings.ToLower(strings.TrimSpace(strFacebook))
}

// CleanEmail returns email as lower case trimmed string.
func CleanEmail(strEmail string) string {
	return strings.ToLower(strings.TrimSpace(strEmail))
}

// CleanLowerCase returns lower case trimmed string.
func CleanLowerCase(strTabLayout string) string {
	return strings.ToLower(strings.TrimSpace(strTabLayout))
}

// CleanUsername returns a username with problem characters cleaned out
func CleanUsername(strUsername string) string {
	strTrim := strings.TrimSpace(strUsername)
	runes := make([]string, 0, 30)
	buf := make([]byte, 4) // bucket for one rune
	var n int
	legalPunct := ` .,_+-&|()?/#` // space char is included

	for _, c := range strTrim {
		switch {
		case strings.ContainsRune(legalPunct, c):
			// include legal punctuation
			n = utf8.EncodeRune(buf, c)
			runes = append(runes, string(buf[:n]))
		case unicode.IsPunct(c):
			// exclude other punctuation
		case unicode.IsSpace(c):
			// exclude other whitespace
		case c == '@':
			// exclude @
		default:
			n = utf8.EncodeRune(buf, c)
			runes = append(runes, string(buf[:n])) // append just the rune
		}
	}
	strCleanUsername := strings.Join(runes, "")
	return strCleanUsername
}

// CleanPhone returns a phone number with problem characters cleaned out and '+' optionally inserted
func CleanPhone(strPhone string, bPlus bool) string {
	runes := make([]string, 0, 30)
	buf := make([]byte, 4) // bucket for one rune
	var n int

	if bPlus == true {
		runes = append(runes, "+")
	}
	for _, c := range strPhone {
		switch {
		case unicode.IsDigit(c):
			n = utf8.EncodeRune(buf, c)
			runes = append(runes, string(buf[:n])) // append just the rune
		default:
			// exclude everything else
		}
	}
	strCleanPhone := strings.Join(runes, "")
	return strCleanPhone
}

// FancyPhone returns phone in format (999)999-9999
func FancyPhone(strPhone string) string {
	s := CleanPhone(strPhone, false)
	if strings.HasPrefix(s, "1") {
		s = s[1:]
	}

	if s == "" {
		return s
	}
	return fmt.Sprintf("(%s)%s-%s", s[:3], s[3:6], s[6:])
}

// ShortURL returns url with http://www. or http: stripped off
func ShortURL(strURL string) string {
	s := strings.ToLower(strings.TrimSpace(strURL))
	s = strings.Replace(strURL, " ", "", -1)

	switch {
	case strings.HasPrefix(s, "http://www."):
		s = s[len("http://www."):]
	case strings.HasPrefix(s, "http://"):
		s = s[len("http://"):]
	case strings.HasPrefix(s, "www."):
		s = s[len("www."):]
	}
	return s
}

// PaypalURL parses a paypal url and makes corrections
func PaypalURL(strRaw string) (string, bool) {
	strRaw = strings.TrimSpace(strRaw)
	if strRaw == "" {
		return "", false
	}

	Trace.Println("PaypalURL raw  ", strRaw)
	u, err := url.Parse(strRaw)
	if err != nil {
		Warning.Println("Bad Paypal url", strRaw, err.Error())
		return strRaw, false
	}

	if u.Scheme == "" {
		Trace.Println("PaypalURL needs scheme added")
		u.Scheme = "http"
		strRaw = u.String()
		u, err = url.Parse(strRaw)
		if err != nil {
			Warning.Println("Bad Paypal url", strRaw, err.Error())
			return strRaw, false
		}
	}

	u.Scheme = "http"
	u.Host = "paypal.me"
	Trace.Println("PaypalURL after ", u.String())
	return u.String(), true
}

// CleanPin returns a pin number with problem characters cleaned out
func CleanPin(strPin string) string {
	runes := make([]string, 0, 30)
	buf := make([]byte, 4) // bucket for one rune
	var n int

	for _, c := range strPin {
		switch {
		case unicode.IsDigit(c):
			n = utf8.EncodeRune(buf, c)
			runes = append(runes, string(buf[:n])) // append just the rune
		default:
			// exclude everything else
		}
	}
	strCleanPin := strings.Join(runes, "")
	return strCleanPin
}

// Reverse reverses order of runes in a string
func Reverse(value string) string {
	data := []rune(value) // not a byte array
	result := []rune{}

	for i := len(data) - 1; i >= 0; i-- {
		result = append(result, data[i])
	}

	return string(result)
}

// CleanPassword returns a password with problem characters cleaned out
func CleanPassword(strPassword string) string {
	runes := make([]string, 0, 30)
	buf := make([]byte, 4) // bucket for one rune
	var n int
	legalPunct := `!-_`

	for _, c := range strPassword {
		switch {
		case strings.ContainsRune(legalPunct, c):
			// include legal punctuation
			n = utf8.EncodeRune(buf, c)
			runes = append(runes, string(buf[:n]))
		case unicode.IsPunct(c):
			// exclude other punctuation
		case unicode.IsSpace(c):
			// exclude whitespace
		case c == '@':
			// exclude @
		default:
			n = utf8.EncodeRune(buf, c)
			runes = append(runes, string(buf[:n])) // append just the rune
		}
	}
	strCleanUsername := strings.Join(runes, "")
	return strCleanUsername
}

// ParseAPIVersion returns api version if found as prefix to url path, else 0
func ParseAPIVersion(strURLPath string) int {
	var intAPIVersion int

	pathParts := strings.SplitN(strURLPath, "/", 3)
	if len(pathParts) > 1 {
		strAPIVersion := CleanAppString(pathParts[1])
		if strings.HasPrefix(strAPIVersion, "v") {
			strNumber := strAPIVersion[1:]
			intAPIVersion, _ = strconv.Atoi(strNumber)
		}
	}
	return intAPIVersion
}

// ParseAPIVersionString returns api version as string if found, else 0
func ParseAPIVersionString(strURLPath string) string {
	var strAPIVersion string

	pathParts := strings.SplitN(strURLPath, "/", 3)
	if len(pathParts) > 1 {
		strAPIVersion = pathParts[1]
		if strings.HasPrefix(strAPIVersion, "v") {
			strAPIVersion = strAPIVersion[1:]
		}
	}
	return strAPIVersion
}

// ParseDraft returns 1 if string contains "/draft/", else 0
func ParseDraft(strURLPath string) int {
	if strings.Contains(strURLPath, "/draft/") {
		return 1
	}
	return 0
}

// ParseAppString parses off the app identifying string from the user agent
// and forces lower case alphanumeric
func ParseAppString(strUserAgent string) string {
	agentParts := strings.SplitN(strUserAgent, "/", 2)
	return CleanAppString(agentParts[0])
}

// CleanAppString forces lower case alphanumeric
func CleanAppString(strApp string) string {
	runes := make([]string, 0, 30)
	buf := make([]byte, 4) // bucket for one rune
	var n int

	for _, c := range strApp {
		switch {
		case unicode.IsDigit(c):
			n = utf8.EncodeRune(buf, c)
			runes = append(runes, string(buf[:n])) // append just the rune
		case unicode.IsLetter(c):
			c = unicode.ToLower(c)
			n = utf8.EncodeRune(buf, c)
			runes = append(runes, string(buf[:n])) // append just the rune
		default:
			// exclude everything else
		}
	}
	strCleanApp := strings.Join(runes, "")
	return strCleanApp
}

// ZeroPad returns a string of the desired length padded with 0's
func ZeroPad(strNumber string, intLen int) string {
	strZeros := "0000000000" + strNumber
	intOffset := len(strZeros) - intLen
	return strZeros[intOffset:]
}

// IntParamAsBool returns optional integer parameter as a boolean
func IntParamAsBool(strParam string) bool {
	intParam, _ := strconv.Atoi(strParam)
	return (intParam != 0)
}

// IntAsBool returns integer as a boolean
func IntAsBool(intParam int) bool {
	return (intParam != 0)
}

// ParamAsInt returns optional integer parameter as int
func ParamAsInt(strParam string) int {
	intParam, _ := strconv.Atoi(strParam)
	return intParam
}

// ParamAsFloat64 returns optional float64 parameter as float64
func ParamAsFloat64(strParam string) float64 {
	dblParam, _ := strconv.ParseFloat(strParam, 64)
	return dblParam
}

// IsNilInterface returns true if either the type is nil
// or content of an interface is empty.
func IsNilInterface(i interface{}) bool {
	if i == nil {
		return true
	}
	t := reflect.TypeOf(i)
	if reflect.ValueOf(i) == reflect.Zero(t) {
		return true
	}
	return false
}
