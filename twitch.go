package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Package level settings.
var baseURI string
var httpClient *http.Client

func init() {
	baseURI = "https://api.twitch.tv/kraken"
	httpClient = http.DefaultClient
}

// A Client contains all of the fields necessary for making requests to the Twitch API.
type Client struct {
	clientID    string
	secret      string
	redirectURI string
}

// An Access holds an access token along with the authorization scope associated.
type Access struct {
	Token string   `json:"access_token"`
	Scope []string `json:"scope"`
}

// An AccessClient wraps a twitch client with an Access struct.
type AccessClient struct {
	access Access
	client Client
}

// NewClient creates a new Client for communicating with the Twitch API.
func NewClient(clientID, secret, redirectURI string) Client {
	return Client{
		clientID:    clientID,
		secret:      secret,
		redirectURI: redirectURI,
	}
}

// NewAccess creates an Access struct from an existing token/scope combination.
func NewAccess(token string, scope []string) Access {
	return Access{
		Token: token,
		Scope: scope,
	}
}

func (c Client) ClientID() string {
	return c.clientID
}

func (c Client) Secret() string {
	return c.secret
}

func (c Client) RedirectURI() string {
	return c.redirectURI
}

// WithAccess wraps the Client with an Access struct.
func (c Client) WithAccess(access Access) AccessClient {
	return AccessClient{
		access: access,
		client: c,
	}
}

// Authorize is an http handler that can be used to prompt a user for Authorization.
func (c Client) Authorize(scope ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authURI := c.getAuthorizeURI(scope)

		http.Redirect(w, r, authURI.String(), 302)
	}
}

// HandleAuthorization handles collecting auth information after a successful auth attempt with Twitch. Once this handler is called, the Client struct should contain the user auth token and scope.
func (c Client) HandleAuthorization(handleAccess func(access Access, err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		handleAccess(c.getAccessToken(query.Get("code")))
	}
}

func (c Client) getAccessToken(authCode string) (Access, error) {
	var access Access
	uri, _ := url.Parse(fmt.Sprintf("%s/oauth2/token", baseURI))
	query := url.Values{}

	query.Add("client_id", c.ClientID())
	query.Add("client_secret", c.Secret())
	query.Add("code", authCode)
	query.Add("grant_type", "authorization_code")
	query.Add("redirect_uri", c.RedirectURI())

	uri.RawQuery = query.Encode()

	req, err := http.NewRequest("POST", uri.String(), nil)
	if err != nil {
		return access, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return access, err
	}

	if err := json.NewDecoder(res.Body).Decode(&access); err != nil {
		return access, err
	}

	return access, nil
}

func (c Client) getAuthorizeURI(scope []string) *url.URL {
	uri, _ := url.Parse(fmt.Sprintf("%s/oauth2/authorize", baseURI))
	queryString := url.Values{}

	queryString.Add("client_id", c.ClientID())
	queryString.Add("redirect_uri", c.RedirectURI())
	queryString.Add("response_type", "code")
	queryString.Add("scope", strings.Join(scope, " "))

	uri.RawQuery = queryString.Encode()
	return uri
}

func (c Client) makeGetRequest(uri string) (*http.Response, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID())
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")

	return httpClient.Do(req)
}

func (ac AccessClient) makeGetRequest(uri string) (*http.Response, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", ac.access.Token))
	req.Header.Set("Client-ID", ac.client.ClientID())
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")

	return httpClient.Do(req)
}

func (c Client) makePostRequest(uri string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", c.ClientID())
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")

	return httpClient.Do(req)
}

func (ac AccessClient) makePostRequest(uri string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", ac.access.Token))
	req.Header.Set("Client-ID", ac.client.ClientID())
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Content-Type", "application/json")

	return httpClient.Do(req)
}

func (ac AccessClient) makePutRequest(uri string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", ac.access.Token))
	req.Header.Set("Client-ID", ac.client.ClientID())
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Content-Type", "application/json")

	return httpClient.Do(req)
}

func (ac AccessClient) makeDeleteRequest(uri string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", ac.access.Token))
	req.Header.Set("Client-ID", ac.client.ClientID())
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")

	return httpClient.Do(req)
}

func (a Access) ValidateScope(scope string) error {
	for _, s := range a.Scope {
		if s == scope {
			return nil
		}
	}

	return fmt.Errorf("Can not complete request because Access struct does not have '%s' scope", scope)
}

// helper function to call validateScope directly on the AccessClient.
func (ac AccessClient) validateScope(scope string) error {
	return ac.access.ValidateScope(scope)
}
