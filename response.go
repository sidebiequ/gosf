package gosf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	errResponse struct {
		status     string
		statusCode int
		errors     []*sfErr
	}

	sfErr struct {
		Message   string `json:"message"`
		ErrorCode string `json:"errorCode"`
	}
)

func (e *sfErr) Error() string {
	return e.ErrorCode + ":" + e.Message
}

func (e *errResponse) Error() string {
	errStr := fmt.Sprintf("error response(%d %s) with %d errors:\n", e.statusCode, e.status, len(e.errors))
	for i, err := range e.errors {
		errStr += fmt.Sprintln(i, ":", err)
	}
	return errStr
}

// parseErrResponse parses the error response from salesforce(response.StatusCode/100!=2)
func parseErrResponse(resp *http.Response) (err error) {
	if resp.StatusCode/100 == 2 {
		err = fmt.Errorf("reponse with success code %d can not be treated as an error response", resp.StatusCode)
		return
	}

	errResp := &errResponse{
		status:     resp.Status,
		statusCode: resp.StatusCode,
	}
	if err = json.NewDecoder(resp.Body).Decode(&errResp.errors); err != nil {
		return
	}

	err = errResp
	return
}
