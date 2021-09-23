package qencode

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetTokenResponse is the get token endpoint response.
type GetTokenResponse struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

// UnmarshalJSON helps to unmarshal time.Time.
func (g *GetTokenResponse) UnmarshalJSON(b []byte) error {
	type Alias GetTokenResponse
	var gt struct {
		Alias
		Expire string `json:"expire"`
	}
	err := json.Unmarshal(b, &gt)
	if err != nil {
		return err
	}

	*g = GetTokenResponse(gt.Alias)
	g.Expire, err = time.Parse("2006-01-02T15:04:05", gt.Expire)
	return err
}

// GetToken retrieves an access token based on the
// api key.
func (c *Client) GetToken(apiKey string) (GetTokenResponse, error) {
	tokenEndpoint := c.baseURL + "/v1/access_token"

	bd := url.Values{}
	bd.Add("api_key", apiKey)
	r, err := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(bd.Encode()))
	if err != nil {
		return GetTokenResponse{}, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.hClient.Do(r)
	if err != nil {
		return GetTokenResponse{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode > 299 {
		b := new(bytes.Buffer)
		_, _ = io.Copy(b, resp.Body)
		return GetTokenResponse{}, RequestError{
			Message:      "error getting qEncode token",
			ResponseBody: b,
			StatusCode:   resp.StatusCode,
		}
	}

	var getTokenResp GetTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&getTokenResp); err != nil {
		return GetTokenResponse{}, err
	}

	return getTokenResp, nil
}
