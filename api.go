package gosf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// follow the Integration: REST API Cheat Sheet	to develop:
// http://resources.docs.salesforce.com/rel1/doc/en-us/static/pdf/SF_Rest_API_cheatsheet_web.pdf

/************************************/
/********* SOBJECT RESOURCES ********/
/************************************/

// CreaetSobject creates sobject by given sobject name and entity.
func (c *Client) CreaetSobject(sobjectName string, sobject interface{}) (string, error) {
	op := &opCreate{
		sobjectName: sobjectName,
		sobject:     sobject,
	}
	if err := c.do(op); err != nil {
		return "", err
	}
	return op.result.ID, nil
}

// UpdateSobject updates sobject by given sobject name, id and entity contains changes.
func (c *Client) UpdateSobject(sobjectName, sobjectID string, sobject interface{}) error {
	return c.do(&opUpdate{
		sobjectName: sobjectName,
		sobjectID:   sobjectID,
		sobject:     sobject,
	})
}

// DeleteSobject deletes sobject by given sobject name and id.
func (c *Client) DeleteSobject(sobjectName, sobjectID string) error {
	return c.do(&opDelete{
		sobjectName: sobjectName,
		sobjectID:   sobjectID,
	})
}

// GetSobject get sobject by given sobject name and id, use target to receive the result.
func (c *Client) GetSobject(sobjectName, sobjectID string, target interface{}) error {
	op := &opGet{
		sobjectName: sobjectName,
		sobjectID:   sobjectID,
	}
	if err := c.do(op); err != nil {
		return err
	}

	byts, _ := json.Marshal(op.target)
	return json.Unmarshal(byts, &target)
}

// QuerySobject query sobject or sobjects by given op OpQuery.
// See also OpQuery.
func (c *Client) QuerySobject(op *OpQuery, target interface{}) error {
	if err := c.do(op); err != nil {
		return err
	}
	byts, _ := json.Marshal(op.target)
	return json.Unmarshal(byts, &target)
}

/************************************/
/****** OTHER COMMON RESOURCES ******/
/************************************/

// Version type. An array with elements describe the versions of salesforce rest api.
type Version struct {
	Label   string `json:"label"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

// Versions shows all availabel ssalesforce rest api versions.
func (c *Client) Versions() (versions []*Version, err error) {
	req, err := http.NewRequest(http.MethodGet, c.requestCtx.BaseURL(), nil)
	if err != nil {
		return
	}

	err = c.doWithHTTPRequest(req, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("verison api can't handle response with code %d, expect %d", resp.StatusCode, http.StatusOK)
		}
		return json.NewDecoder(resp.Body).Decode(&versions)
	})
	return
}

// Resources shows resources under current version.
func (c *Client) Resources() (resources map[string]string, err error) {
	req, err := http.NewRequest(http.MethodGet, c.requestCtx.VersionURL(), nil)
	if err != nil {
		return
	}

	err = c.doWithHTTPRequest(req, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("resources api can't handle response with code %d, expect %d", resp.StatusCode, http.StatusOK)
		}
		return json.NewDecoder(resp.Body).Decode(&resources)
	})
	return
}

// SobjectInfo shows the basic information of given sobject name.
func (c *Client) SobjectInfo() (info map[string]interface{}, err error) {
	req, err := http.NewRequest(http.MethodGet, c.requestCtx.SobjectURL(), nil)
	if err != nil {
		return
	}

	err = c.doWithHTTPRequest(req, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("resources api can't handle response with code %d, expect %d", resp.StatusCode, http.StatusOK)
		}
		return json.NewDecoder(resp.Body).Decode(&info)
	})
	return
}
