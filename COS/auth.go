package COS

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Timeout        = 1 * time.Hour
	qSignAlgorithm = "sha1"
)

//implement the cos sign procedure
func sign(req *http.Request, screctKey, secretID string) {
	qSignAlgorithm := qSignAlgorithm
	qAk := secretID
	qSignTime := signTimeoutTime()
	qKeyTime := qSignTime
	qHeaderList := strings.Join(headerList(req), ";")
	qURLParamList := strings.Join(paramList(req), ";")
	qSignature := signature(req, screctKey)
	authHeader := fmt.Sprintf("q-sign-algorithm=%s&q-ak=%s&q-sign-time=%s&q-key-time=%s&q-header-list=%s&q-url-param-list=%s&q-signature=%s",
		qSignAlgorithm, qAk, qSignTime, qKeyTime, qHeaderList, qURLParamList, qSignature)
	req.Header.Add("Authorization", authHeader)
	fmt.Println(authHeader)
	fmt.Println("headers", req.Header)
}

//q-sign-time, q-key-time
func signTimeoutTime() string {
	now := time.Now()
	end := now.Add(Timeout)
	res := strconv.Itoa(int(now.Unix())) + ";" + strconv.Itoa(int(end.Unix()))
	return res
}

//q-header-list
func headerList(req *http.Request) []string {
	var list []string
	for key := range req.Header {
		list = append(list, strings.ToLower(key))
	}
	sort.Strings(list)
	return list
}

//q-url-param-list
func paramList(req *http.Request) []string {
	var list []string
	v := req.URL.Query()
	for key := range v {
		list = append(list, strings.ToLower(key))
	}
	sort.Strings(list)
	return list
}

//sha1 hmac algorithm to sign the key
func signKey(key, qKeyTime string) string {
	h := hmac.New(sha1.New, []byte(key))
	io.WriteString(h, qKeyTime)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func httpString(req *http.Request) string {
	method := strings.ToLower(req.Method)
	uri := req.URL.Path //EscapedPath??
	params := httpParams(req)
	headers := httpHeaders(req)
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", method, uri, params, headers)
}

//HttpParameters
func httpParams(req *http.Request) string {
	list := paramList(req)
	params := req.URL.Query()
	params = toLower(params)
	res := ""
	for _, k := range list {
		if res != "" {
			res += "&"
		}
		v := params.Get(k)
		res += k + "=" + v
	}
	return res
}

//convert url.Values key to lower case
func toLower(m url.Values) url.Values {
	res := make(url.Values)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}
	return res
}

func httpHeaders(req *http.Request) string {
	list := headerList(req)
	var res string
	for _, k := range list {
		if res != "" {
			res += "&"
		}
		v := req.Header.Get(k)
		res += k + "=" + url.PathEscape(v)
	}
	return res
}

func stringToSign(req *http.Request, qSignTime string) string {
	HTTPString := httpString(req)
	h := sha1.New()
	io.WriteString(h, HTTPString)
	return fmt.Sprintf("%s\n%s\n%s\n", qSignAlgorithm, qSignTime, h.Sum(nil))
}

func signature(req *http.Request, SecretKey string) string {
	qSignTime := signTimeoutTime()
	signKey := signKey(SecretKey, qSignTime)
	h := hmac.New(sha1.New, []byte(signKey))
	io.WriteString(h, stringToSign(req, qSignTime))
	return fmt.Sprintf("%x", h.Sum(nil))
}
