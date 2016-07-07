package GoSalesforce

import (
	"fmt"
	"net/http"
	"strings"
)

/*************************************/
/************** REQUEST **************/
/*************************************/

// Request interface defines what a request for salesforce restful api should provide:
//  - the request method.
//  - the request url generator which can generate url by given host and api version.
//  - the data can be jsonify and will be send.
//  - the http code when operation is success.
//  - the method to verify whether this request is valid.
type Request interface {
	// Method defines what http method should the request use.
	Method() string

	// Method defines what http method should the request use.
	URL(ctx *RequestCtx) string

	// URL returns a URLGenerator to generate url by given host and api version.
	Data() interface{}

	// SuccessCode defines the response http code when request success.
	SuccessCode() int

	// IsValid returns true if the request's field are all valid.
	IsValid() bool
}

var (
	_ Request = &RequestCreateSObject{}
	_ Request = &RequestUpdateSObject{}
	_ Request = &RequestDeleteSObject{}
	_ Request = &RequestQuerySObject{}
)

/************************************/
/********** CREATE SOBJECT **********/
/************************************/

// RequestCreateSObject is a request for creating SObject.
type RequestCreateSObject struct {
	SObjectName string
	SObject     interface{}
}

// Method defines what http method should the request use.
func (r *RequestCreateSObject) Method() string {
	return http.MethodPost
}

// URL returns a url to be sent by given request context.
func (r *RequestCreateSObject) URL(ctx *RequestCtx) string {
	return ctx.SObjectURL(r.SObjectName)
}

// Data returns the request data or nil if the request send nothing.
func (r *RequestCreateSObject) Data() interface{} {
	return r.Data
}

// SuccessCode defines the response http code when request success.
func (r *RequestCreateSObject) SuccessCode() int {
	return http.StatusCreated
}

// IsValid returns true if the request's field are all valid.
func (r *RequestCreateSObject) IsValid() bool {
	return r.SObjectName != "" && r.SObject != nil
}

/************************************/
/********** UPDATE SOBJECT **********/
/************************************/

// RequestUpdateSObject is a request for creating SObject.
type RequestUpdateSObject struct {
	SObjectName string
	SObjectID   string
	SObject     interface{}
}

// Method defines what http method should the request use.
func (r *RequestUpdateSObject) Method() string {
	return http.MethodPatch
}

// URL returns a url to be sent by given request context.
func (r *RequestUpdateSObject) URL(ctx *RequestCtx) string {
	return ctx.SObjectOpURL(r.SObjectName, r.SObjectID)
}

// Data returns the request data or nil if the request send nothing.
func (r *RequestUpdateSObject) Data() interface{} {
	return r.Data
}

// SuccessCode defines the response http code when request success.
func (r *RequestUpdateSObject) SuccessCode() int {
	return http.StatusOK
}

// IsValid returns true if the request's field are all valid.
func (r *RequestUpdateSObject) IsValid() bool {
	return r.SObjectName != "" && r.SObjectID != "" && r.SObject != nil
}

/************************************/
/********** DELETE SOBJECT **********/
/************************************/

// RequestDeleteSObject is a request for creating SObject.
type RequestDeleteSObject struct {
	SObjectName string
	SObjectID   string
}

// Method defines what http method should the request use.
func (r *RequestDeleteSObject) Method() string {
	return http.MethodDelete
}

// URL returns a url to be sent by given request context.
func (r *RequestDeleteSObject) URL(ctx *RequestCtx) string {
	return ctx.SObjectOpURL(r.SObjectName, r.SObjectID)
}

// Data returns the request data or nil if the request send nothing.
func (r *RequestDeleteSObject) Data() interface{} {
	return nil
}

// SuccessCode defines the response http code when request success.
func (r *RequestDeleteSObject) SuccessCode() int {
	return http.StatusNoContent
}

// IsValid returns true if the request's field are all valid.
func (r *RequestDeleteSObject) IsValid() bool {
	return r.SObjectName != "" && r.SObjectID != ""
}

/***********************************/
/********** QUERY SOBJECT **********/
/***********************************/

type whereClause struct {
	field     string
	condition interface{}
}

// IsValid returns true if whereClause's condition is valid.
// In SOQL, condition in where clause can only be number, boolean or string.
func (c *whereClause) IsValid() bool {
	switch c.condition.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case bool:
		return true
	case string:
		return true
	default:
		return false
	}
}

// RequestQuerySObject is a request for quering SObject.
type RequestQuerySObject struct {
	sObjectName  string
	selectFileds []string
	whereClauses []whereClause
	order        string
	nullPriority string
	limit        int
}

// Method defines what http method should the request use.
func (r *RequestQuerySObject) Method() string {
	return http.MethodGet
}

// URL returns a url to be sent by given request context.
func (r *RequestQuerySObject) URL(ctx *RequestCtx) string {
	query := r.makeQueryStatment()
	return ctx.QueryURL(query)
}

// Data returns the request data or nil if the request send nothing.
func (r *RequestQuerySObject) Data() interface{} {
	return nil
}

// SuccessCode defines the response http code when request success.
func (r *RequestQuerySObject) SuccessCode() int {
	return http.StatusOK
}

// IsValid returns true if the request's field are all valid.
func (r *RequestQuerySObject) IsValid() bool {
	return len(r.selectFileds) > 0 && r.sObjectName != ""
}

// Select defines which columns of SObject will be return in result.
func (r *RequestQuerySObject) Select(fields ...string) *RequestQuerySObject {
	if r.selectFileds == nil {
		r.selectFileds = make([]string, 0)
	}
	r.selectFileds = append(r.selectFileds, fields...)
	return r
}

// From defines which SObject will be query.
func (r *RequestQuerySObject) From(sObjectName string) *RequestQuerySObject {
	r.sObjectName = sObjectName
	return r
}

// OrderDesc defines the order column to sort results descendant.
func (r *RequestQuerySObject) OrderDesc(field string) *RequestQuerySObject {
	r.order = fmt.Sprintf("%s DESC", field)
	return r
}

// OrderAsc defines the order column to sort results ascendant.
func (r *RequestQuerySObject) OrderAsc(field string) *RequestQuerySObject {
	r.order = fmt.Sprintf("%s ASC", field)
	return r
}

// OrderReset resets the order statment=="".
func (r *RequestQuerySObject) OrderReset() *RequestQuerySObject {
	r.order = ""
	return r
}

// OrderNullFirst make null null column value first in query results when given order column.
func (r *RequestQuerySObject) OrderNullFirst() *RequestQuerySObject {
	r.nullPriority = "NULL FIRST"
	return r
}

// OrderNullLast make null column value last in query results when given order column.
func (r *RequestQuerySObject) OrderNullLast() *RequestQuerySObject {
	r.nullPriority = "NULL LAST"
	return r
}

// Limit defines the max records will be return.
// if n==0, will treat it as a signal to reset limitation to unlimited.
func (r *RequestQuerySObject) Limit(n int) *RequestQuerySObject {
	r.limit = n
	return r
}

func (r *RequestQuerySObject) makeQueryStatment() string {
	baseQuery := fmt.Sprintf("SELECT %s FROM %s", r.makeSelectStatment(), r.sObjectName)
	return strings.Join([]string{
		baseQuery,
		r.makeWhereCluasesStatment(),
		r.makeOrderStatment(),
		r.makeLimitStatment(),
	}, " ")
}

// makeSelectStatment renders statment as below if r.selectFileds has elements:
// SELECT <FIELD1> [,<FIELD2>]...
func (r *RequestQuerySObject) makeSelectStatment() string {
	if len(r.selectFileds) == 0 {
		defaultLogger.Print("[QuerySObjectRequest] Missing Select fields")
	}
	return strings.Join(r.selectFileds, ",")
}

// makeWhereCluasesStatment renders statment as below if r.whereClauses has elements:
// WHERE <FIELD1>=<CODITION1> [,<FIELD2>=<CODITION2>]...
func (r *RequestQuerySObject) makeWhereCluasesStatment() string {
	if len(r.whereClauses) == 0 {
		return ""
	}

	var data = make(map[string]interface{})
	for _, clause := range r.whereClauses {
		if !clause.IsValid() {
			defaultLogger.Printf(
				"[QuerySObjectRequest] Found invalid where-clause which should be type of number, boolean or string. QuerySObject=%s, Field=%s, ConditionType=%T",
				r.sObjectName, clause.field, clause.condition,
			)
			continue
		}
		data[clause.field] = clause.condition
	}

	var filters = make([]string, 0)
	for field, cond := range data {
		filters = append(filters, fmt.Sprintf("%s=%v", field, cond))
	}
	return fmt.Sprintf("WHERE %s", strings.Join(filters, ","))
}

// makeOrderStatment renders statment as below if r.order!="":
// ORDER BY <FIELD> <DESC|ASC> NULL <FIRST|LAST>
func (r *RequestQuerySObject) makeOrderStatment() string {
	if r.order == "" {
		return ""
	}

	var nullPriority = r.nullPriority
	if r.nullPriority == "" {
		nullPriority = "NULL FIRST"
	}
	return fmt.Sprintf("ORDER BY %s %s", r.order, nullPriority)
}

// makeLimitStatment renders statment as below if r.limit>0:
// LIMIT <LIMITATION>
func (r *RequestQuerySObject) makeLimitStatment() string {
	if r.limit > 0 {
		return fmt.Sprintf("LIMIT %d", r.limit)
	}
	return ""
}
