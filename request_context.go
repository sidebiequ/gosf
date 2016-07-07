package GoSalesforce

import "fmt"
import "net/url"

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

// SObjectBaseURL returns the URL can work with SObjects, like:
// "https://instance.salesforce.com/services/data/v36.0/sobjects"
func (ctx *RequestCtx) SObjectBaseURL() string {
	return fmt.Sprintf("%s/sobjects", ctx.VersionURL())
}

// SObjectURL returns the URL with specific sobject.
// Assume the given sobject is 'User', the return URL will be:
// "https://instance.salesforce.com/services/data/v36.0/sobjects/User"
func (ctx *RequestCtx) SObjectURL(sobject string) string {
	return fmt.Sprintf("%s/%s", ctx.SObjectBaseURL(), sobject)
}

// SObjectOpURL returns the URL with specific sobject and id.
// Assume the given sobject is 'User' and id is '00e28000001K04LAA1',
// the return URL will be:
// "https://instance.salesforce.com/services/data/v36.0/sobjects/User/00e28000001K04LAA1"
func (ctx *RequestCtx) SObjectOpURL(sobject, sobjectID string) string {
	return fmt.Sprintf("%s/%s", ctx.SObjectURL(sobject), sobjectID)
}

func (ctx *RequestCtx) isVersionValid() bool {
	return ctx.version >= minAPIVersion && ctx.version <= maxAPIVersion
}
