package gosf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

/*************************************/
/************** REQUEST **************/
/*************************************/

// Request can make a http.Request instance which contains json-format body if data provided.
type Request struct {
	method string
	urlStr string
	data   interface{}
}

// NewRequest returns a new Request given a method, URL, and optional data.
func NewRequest(method, urlStr string, data interface{}) *Request {
	return &Request{
		method: method,
		urlStr: urlStr,
		data:   data,
	}
}

func (r *Request) makeRequest() (*http.Request, error) {
	if r.method == "" {
		return nil, errors.New("missing request method")
	}

	if _, err := url.Parse(r.urlStr); err != nil {
		return nil, err
	}
	fmt.Println("r is %s", r.data)
	if r.data == nil {
		return http.NewRequest(r.method, r.urlStr, nil)
	}
	fmt.Println("r is %s", r.data)
	byts, err := json.Marshal(r.data)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(r.method, r.urlStr, bytes.NewReader(byts))

}

/*************************************/
/********** REQUEST CONTEXT **********/
/*************************************/

const (
	minAPIVersion = 8
	maxAPIVersion = 37
)

// RequestCtx holds the host and api version informations.
type RequestCtx struct {
	host    string
	version int
}

// BaseURL returns the base URL of salesforce restful api.
// Assume host is "https://instance.salesforce.com", baseURL will be
// "https://instance.salesforce.com/services/data"
func (ctx *RequestCtx) BaseURL() string {
	return fmt.Sprintf("%s/services/data", ctx.host)
}

// VersionURL returns the URL with version, like:
// "https://instance.salesforce.com/services/data/v36.0"
func (ctx *RequestCtx) VersionURL() string {
	return fmt.Sprintf("%s/v%d.0", ctx.BaseURL(), ctx.version)
}

// QueryURL returns the URL with query SOQL statments, like:
// "https://instance.salesforce.com/services/data/v24.0/query?q=SELECT+Id,+Name+FROM+User"
func (ctx *RequestCtx) QueryURL(q string) string {
	return fmt.Sprintf("%s/query?q=%s", ctx.VersionURL(), url.QueryEscape(q))
}

// SobjectURL returns the URL can work with SObjects, like:
// "https://instance.salesforce.com/services/data/v36.0/sobjects"
func (ctx *RequestCtx) SobjectURL() string {
	return fmt.Sprintf("%s/sobjects", ctx.VersionURL())
}

// SobjectURLWithName returns the URL with specific sobject.
// Assume the given sobject is 'User', the return URL will be:
// "https://instance.salesforce.com/services/data/v36.0/sobjects/User"
func (ctx *RequestCtx) SobjectURLWithName(sobjectName string) string {
	return fmt.Sprintf("%s/%s", ctx.SobjectURL(), sobjectName)
}

// SobjectURLWithID returns the URL with specific sobject and id.
// Assume the given sobject is 'User' and id is '00e28000001K04LAA1',
// the return URL will be:
// "https://instance.salesforce.com/services/data/v36.0/sobjects/User/00e28000001K04LAA1"
func (ctx *RequestCtx) SobjectURLWithID(sobjectName, sobjectID string) string {
	return fmt.Sprintf("%s/%s", ctx.SobjectURLWithName(sobjectName), sobjectID)
}

func (ctx *RequestCtx) isVersionValid() bool {
	return ctx.version >= minAPIVersion && ctx.version <= maxAPIVersion
}
