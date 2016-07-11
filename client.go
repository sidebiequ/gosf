package gosf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

/************************************/
/************** CONFIG **************/
/************************************/

// Config type
type Config struct {
	// auth
	Host         string `json:"host"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ExpiresIn    int    `json:"expires_in"`

	// api
	APIVersion int `json:"api_version"`
}

/***********************************/
/************** TOKEN **************/
/***********************************/

// Token type.
type token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	Signature   string    `json:"signature"`
	ExpiresAt   time.Time `json:"-"`
}

// IsExpired returns true if token has expired.
func (t *token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

/***********************************/
/************** OAUTH **************/
/***********************************/

// OAuth type
type oAuth struct {
	*Config
	*token
}

// ExchangeToken exchange token by username and password.
// If fail, salesforce will retrun
//      {
//          "error": "ERROR_TYPE",
//          "error_description": "ERROR_DESCRIPTION"
//      }
func (o *oAuth) exchangeToken() (t *token, err error) {
	u := o.Host + "/services/oauth2/token"
	form := url.Values{
		"grant_type":    {"password"},
		"client_id":     {o.ClientID},
		"client_secret": {o.ClientSecret},
		"username":      {o.Username},
		"password":      {o.Password},
	}

	resp, err := http.PostForm(u, form)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		var authErrResp struct {
			Err            string `json:"error"`
			ErrDescription string `json:"error_description"`
		}
		err = decoder.Decode(&authErrResp)
		if err != nil {
			return
		}

		err = fmt.Errorf("OAuth fail(%s): %s", authErrResp.Err, authErrResp.ErrDescription)
		return
	}

	t = &token{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(o.ExpiresIn)),
	}
	err = decoder.Decode(t)
	if err != nil {
		return
	}
	return
}

func (o *oAuth) RoundTrip(req *http.Request) (res *http.Response, err error) {
	if o.token == nil || o.token.IsExpired() {
		o.token, err = o.exchangeToken()
		if err != nil {
			return
		}
	}
	req.Header.Add("Authorization", o.TokenType+" "+o.AccessToken)
	return http.DefaultTransport.RoundTrip(req)
}

/************************************/
/************** CLIENT **************/
/************************************/

// Client Type
type Client struct {
	client     *http.Client
	requestCtx *RequestCtx
	logger     Logger
}

// NewClient returns a Client instance.
func NewClient(config *Config, logger Logger) *Client {
	// set logger to defaultLogger if the logger argument is nil
	if logger == nil {
		logger = defaultLogger
		logger.Print("[logger] argument 'logger' is nil, use defaultLogger instead")
	} else {
		defaultLogger = logger
	}

	// set config.ExpiresIn default value if it's invalid
	if config.ExpiresIn <= 0 {
		logger.Printf("[expired time] invalid token expired time received, set to 3600s instead")
		config.ExpiresIn = 3600
	}
	logger.Printf("[expired time] token expired after %ds and auto-refresh", config.ExpiresIn)

	// set requestCtx.version to maxVersion if it is out of bound
	requestCtx := &RequestCtx{
		host:    config.Host,
		version: config.APIVersion,
	}
	if !requestCtx.isVersionValid() {
		requestCtx.version = maxAPIVersion
		logger.Printf("[api version] config.APIVersion is out of bound [%d, %d], set to %d instead", minAPIVersion, maxAPIVersion, requestCtx.version)
	}

	return &Client{
		client: &http.Client{
			Transport: &oAuth{Config: config},
		},
		requestCtx: requestCtx,
		logger:     logger,
	}
}

func (c *Client) do(op Operator) (err error) {
	req, err := op.Make(&RequestCtx{
		host:    c.requestCtx.host,
		version: c.requestCtx.version,
	})
	if err != nil {
		return
	}

	httpReq, err := req.makeRequest()
	if err != nil {
		return
	}

	return c.doWithHTTPRequest(httpReq, op.Handle)
}

func (c *Client) doWithHTTPRequest(httpReq *http.Request, handler func(*http.Response) error) (err error) {
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 == 2 {
		err = handler(resp)
	} else {
		err = parseErrResponse(resp)
	}
	return
}

// Do uses op Operator to make request, send it to salesforce,
// receieve response and pass it to the Operator to handle.
func (c *Client) Do(op Operator) error {
	return c.do(op)
}
