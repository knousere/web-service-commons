package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// MakeNonce makes a string to uniquely tag an http request.
func MakeNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", err
	}

	nonce := base64.URLEncoding.EncodeToString(nonceBytes)
	return nonce, nil
}

// MakeKeyValuePairDQ properly formats a key value pair for inclusion in a url.
// The Value is doublequoted.
// RFC 3986 compliant encoding
func MakeKeyValuePairDQ(key string, value string) string {
	encodedKey := url.QueryEscape(key)
	encodedValue := url.QueryEscape(value)
	pair := encodedKey + "=" + DQ(encodedValue)
	return pair
}

// MakeKeyValuePair properly formats a key value pair for inclusion in a signature.
// RFC 3986 compliant encoding
func MakeKeyValuePair(key string, value string) string {
	encodedKey := url.QueryEscape(key)
	encodedValue := url.QueryEscape(value)
	pair := encodedKey + "=" + encodedValue
	return pair
}

// MakeCurrrentTimestamp returns seconds since the beginning of Unix time as a string.
// This is exactly equivalent to the MySQL current_timestamp() function.
func MakeCurrrentTimestamp() string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	return timestamp
}

// MakeSignatureHex returns hex signature string using HMAC-SHA1 hashing
func MakeSignatureHex(strSignatureBase string, strKey string) string {
	key := []byte(strKey)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(strSignatureBase))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MakeSignature64 returns base 64 signature string using HMAC-SHA1 hashing
func MakeSignature64(strSignatureBase string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(strSignatureBase))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// MakeSignatureBase assembles a signature base for input to MakeSignature64()
func MakeSignatureBase(strConsumerKey string, strToken string, strNonce string, strTimestamp string,
	strHTTPMethod string, strURL string, params []KeyValuePair) string {

	sigConsumerKey := MakeKeyValuePair("oauth_consumer_key", strConsumerKey)
	sigNonce := MakeKeyValuePair("oauth_nonce", strNonce)
	sigMethod := MakeKeyValuePair("oauth_signature_method", "HMAC-SHA1")
	sigTimestamp := MakeKeyValuePair("oauth_timestamp", strTimestamp)
	sigToken := MakeKeyValuePair("oauth_token", strToken)
	sigVersion := MakeKeyValuePair("oauth_version", "1.0")

	// create a param string for the signature
	sigParamList := make([]string, 0, 6+len(params))
	sigParamList = append(sigParamList, sigConsumerKey, sigNonce, sigMethod, sigTimestamp, sigToken, sigVersion)

	for _, param := range params {
		sigParam := MakeKeyValuePair(param.Key, param.Value)
		sigParamList = append(sigParamList, sigParam)
	}

	sort.Strings(sigParamList)
	strSigParam := strings.Join(sigParamList, "&")
	Trace.Println("strSigParam ", strSigParam)

	// build the signature
	sigList := make([]string, 0, 3)
	sigList = append(sigList, strHTTPMethod, url.QueryEscape(strURL), url.QueryEscape(strSigParam))
	strSignatureBase := strings.Join(sigList, "&")

	Trace.Println("strSignatureBase ", strSignatureBase)
	return strSignatureBase
}

// MakeHeader constructs an Oauth header string
func MakeHeader(strConsumerKey string, strToken string, strNonce string, strTimestamp string,
	strSignature string) string {

	hdrConsumerKey := MakeKeyValuePairDQ("oauth_consumer_key", strConsumerKey)
	hdrNonce := MakeKeyValuePairDQ("oauth_nonce", strNonce)
	hdrMethod := MakeKeyValuePairDQ("oauth_signature_method", "HMAC-SHA1")
	hdrTimestamp := MakeKeyValuePairDQ("oauth_timestamp", strTimestamp)
	hdrToken := MakeKeyValuePairDQ("oauth_token", strToken)
	hdrVersion := MakeKeyValuePairDQ("oauth_version", "1.0")
	hdrSignature := MakeKeyValuePairDQ("oauth_signature", strSignature)

	// create the authorization header
	hdrParamList := make([]string, 0, 7)
	hdrParamList = append(hdrParamList, hdrConsumerKey, hdrNonce, hdrMethod, hdrSignature, hdrTimestamp, hdrToken, hdrVersion)
	sort.Strings(hdrParamList)
	strHeader := "OAuth " + strings.Join(hdrParamList, ", ")
	return strHeader
}
