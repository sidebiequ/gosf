package gosf

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

/*************************************/
/************** Operator *************/
/*************************************/

// Operator can make a *http.Request instance according to given request ctx
// for http.Client to send request to salesforce or returns an error.
// It can also handle the response responded from salesforce.
type Operator interface {
	// Make request by given request context.
	Make(*RequestCtx) (*Request, error)
	// Handle success response from salesforce.
	Handle(*http.Response) error
}

var (
	_ Operator = &opCreate{}
	_ Operator = &opUpdate{}
	_ Operator = &opDelete{}
	_ Operator = &opGet{}
	_ Operator = &OpQuery{}
)

/************************************/
/********** CREATE SOBJECT **********/
/************************************/

type (
	// opCreate is a request for creating SObject.
	opCreate struct {
		sobjectName string
		sobject     interface{}
		result      *opCreateResult
	}

	opCreateResult struct {
		ID      string   `json:"id"`
		Errors  []string `json:"errors"`
		Success bool     `json:"success"`
	}
)

func (op *opCreate) Make(ctx *RequestCtx) (*Request, error) {
	switch {
	case op.sobjectName == "":
		return nil, errors.New("missing Sobject name")
	case op.sobject == nil:
		return nil, errors.New("missing Sobject")
	default:
		return NewRequest(http.MethodPost, ctx.SobjectURLWithName(op.sobjectName), op.sobject), nil
	}
}

func (op *opCreate) Handle(resp *http.Response) error {
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create operator can't handle response with code %d, expect %d", resp.StatusCode, http.StatusCreated)
	}
	return json.NewDecoder(resp.Body).Decode(&op.result)
}

/************************************/
/********** UPDATE SOBJECT **********/
/************************************/

// opUpdate is a request for updating SObject.
type opUpdate struct {
	sobjectName string
	sobjectID   string
	sobject     interface{}
}

func (op *opUpdate) Make(ctx *RequestCtx) (*Request, error) {
	switch {
	case op.sobjectName == "":
		return nil, errors.New("missing Sobject name")
	case op.sobjectID == "":
		return nil, errors.New("missing Sobject id")
	case op.sobject == nil:
		return nil, errors.New("missing Sobject")
	default:
		return NewRequest(http.MethodPatch, ctx.SobjectURLWithID(op.sobjectName, op.sobjectID), op.sobject), nil
	}
}

func (op *opUpdate) Handle(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update operator can't handle response with code %d, expect %d", resp.StatusCode, http.StatusOK)
	}
	return nil
}

/************************************/
/********** DELETE SOBJECT **********/
/************************************/

// opDelete is a request for deleting SObject.
type opDelete struct {
	sobjectName string
	sobjectID   string
}

func (op *opDelete) Make(ctx *RequestCtx) (*Request, error) {
	switch {
	case op.sobjectName == "":
		return nil, errors.New("missing Sobject name")
	case op.sobjectID == "":
		return nil, errors.New("missing Sobject id")
	default:
		return NewRequest(http.MethodDelete, ctx.SobjectURLWithID(op.sobjectName, op.sobjectID), nil), nil
	}
}

func (op *opDelete) Handle(resp *http.Response) error {
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete operator can't handle response with code %d, expect %d", resp.StatusCode, http.StatusNoContent)
	}
	return nil
}

/************************************/
/*********** GET SOBJECT ************/
/************************************/

// opGet is a request for getting SObject.
type opGet struct {
	sobjectName string
	sobjectID   string
	target      interface{}
}

func (op *opGet) Make(ctx *RequestCtx) (*Request, error) {
	switch {
	case op.sobjectName == "":
		return nil, errors.New("missing Sobject name")
	case op.sobjectID == "":
		return nil, errors.New("missing Sobject id")
	default:
		return NewRequest(http.MethodGet, ctx.SobjectURLWithID(op.sobjectName, op.sobjectID), nil), nil
	}
}

func (op *opGet) Handle(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get operator can't handle response with code %d, expect %d", resp.StatusCode, http.StatusOK)
	}
	return json.NewDecoder(resp.Body).Decode(&op.target)
}

/***********************************/
/********** QUERY SOBJECT **********/
/***********************************/

type (
	// OpQuery is a request for quering SObject.
	// Set which Sobject type and fileds to query by Select() and From() is necessary.
	// The reset operations below is optional:
	// 	- OrderDesc(), OrderAsc(), OrderReset(), OrderNullFirst(), OrderNullLast() to control ORDER key.
	//  - Limit() to control LIMIT key
	// See the methods' doc  for more details.
	OpQuery struct {
		sobjectName  string
		selectFileds []string
		whereClauses []whereClause
		order        string
		nullPriority string
		limit        int
		target       interface{}
	}

	whereClause struct {
		field     string
		condition interface{}
	}
)

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

// Make request by given request context.
func (op *OpQuery) Make(ctx *RequestCtx) (*Request, error) {
	switch {
	case op.sobjectName == "":
		return nil, errors.New("missing Sobject name")
	case len(op.selectFileds) <= 0:
		return nil, errors.New("missing select fields")
	default:
		return NewRequest(http.MethodGet, ctx.QueryURL(op.makeQueryStatment()), nil), nil
	}
}

// Handle success response from salesforce.
func (op *OpQuery) Handle(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get operator can't handle response with code %d, expect %d", resp.StatusCode, http.StatusOK)
	}
	return json.NewDecoder(resp.Body).Decode(&op.target)
}

// Select defines which columns of SObject will be return in result.
func (op *OpQuery) Select(fields ...string) *OpQuery {
	if op.selectFileds == nil {
		op.selectFileds = make([]string, 0)
	}
	op.selectFileds = append(op.selectFileds, fields...)
	return op
}

// From defines which SObject will be query.
func (op *OpQuery) From(sobjectName string) *OpQuery {
	op.sobjectName = sobjectName
	return op
}

// OrderDesc defines the order column to sort results descendant.
func (op *OpQuery) OrderDesc(field string) *OpQuery {
	op.order = fmt.Sprintf("%s DESC", field)
	return op
}

// OrderAsc defines the order column to sort results ascendant.
func (op *OpQuery) OrderAsc(field string) *OpQuery {
	op.order = fmt.Sprintf("%s ASC", field)
	return op
}

// OrderReset resets the order statment=="".
func (op *OpQuery) OrderReset() *OpQuery {
	op.order = ""
	return op
}

// OrderNullFirst make null null column value first in query results when given order column.
func (op *OpQuery) OrderNullFirst() *OpQuery {
	op.nullPriority = "NULL FIRST"
	return op
}

// OrderNullLast make null column value last in query results when given order column.
func (op *OpQuery) OrderNullLast() *OpQuery {
	op.nullPriority = "NULL LAST"
	return op
}

// Limit defines the max records will be return.
// if n==0, will treat it as a signal to reset limitation to unlimited.
func (op *OpQuery) Limit(n int) *OpQuery {
	op.limit = n
	return op
}

func (op *OpQuery) makeQueryStatment() string {
	baseQuery := fmt.Sprintf("SELECT %s FROM %s", op.makeSelectStatment(), op.sobjectName)
	return strings.Join([]string{
		baseQuery,
		op.makeWhereCluasesStatment(),
		op.makeOrderStatment(),
		op.makeLimitStatment(),
	}, " ")
}

// makeSelectStatment renders statment as below if r.selectFileds has elements:
// SELECT <FIELD1> [,<FIELD2>]...
func (op *OpQuery) makeSelectStatment() string {
	if len(op.selectFileds) == 0 {
		defaultLogger.Print("[QuerySObjectRequest] Missing Select fields")
	}
	return strings.Join(op.selectFileds, ",")
}

// makeWhereCluasesStatment renders statment as below if r.whereClauses has elements:
// WHERE <FIELD1>=<CODITION1> [,<FIELD2>=<CODITION2>]...
func (op *OpQuery) makeWhereCluasesStatment() string {
	var data = make(map[string]interface{})
	for _, clause := range op.whereClauses {
		if !clause.IsValid() {
			defaultLogger.Printf(
				"[QuerySObjectRequest] Found invalid where-clause which should be type of number, boolean or string. QuerySObject=%s, Field=%s, ConditionType=%T",
				op.sobjectName, clause.field, clause.condition,
			)
			continue
		}
		data[clause.field] = clause.condition
	}

	if len(data) == 0 {
		return ""
	}

	var filters = make([]string, 0)
	for field, cond := range data {
		filters = append(filters, fmt.Sprintf("%s=%v", field, cond))
	}
	return fmt.Sprintf("WHERE %s", strings.Join(filters, ","))
}

// makeOrderStatment renders statment as below if r.order!="":
// ORDER BY <FIELD> <DESC|ASC> NULL <FIRST|LAST>
func (op *OpQuery) makeOrderStatment() string {
	if op.order == "" {
		return ""
	}

	var nullPriority = op.nullPriority
	if op.nullPriority == "" {
		nullPriority = "NULL FIRST"
	}
	return fmt.Sprintf("ORDER BY %s %s", op.order, nullPriority)
}

// makeLimitStatment renders statment as below if r.limit>0:
// LIMIT <LIMITATION>
func (op *OpQuery) makeLimitStatment() string {
	if op.limit > 0 {
		return fmt.Sprintf("LIMIT %d", op.limit)
	}
	return ""
}

// NewOpQuery returns a OpQuery instance with given sobjectName.
func NewOpQuery(sobjectName string) *OpQuery {
	return &OpQuery{
		sobjectName: sobjectName,
	}
}
