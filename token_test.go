package qencode

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const getTokenResp = ` 
	{
		"token": "1357924680",
		"expire": "2021-09-19T01:35:57"
	}`

func TestClient_GetToken(t *testing.T) {
	apiKey := "api_key"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotAPIKey := r.Form.Get("api_key")
		if gotAPIKey != apiKey {
			t.Fatalf("Client.GetToken() | got api key %s in request, want %s", gotAPIKey, apiKey)
		}
		_, _ = w.Write([]byte(getTokenResp))
	}))
	defer s.Close()

	c := &Client{
		baseURL: s.URL,
		hClient: http.DefaultClient,
	}

	want := GetTokenResponse{
		Token: "1357924680",
		Expire: func() time.Time {
			tm, err := time.Parse("2006-01-02T15:04:05", "2021-09-19T01:35:57")
			if err != nil {
				t.Fatal(err)
			}
			return tm
		}(),
	}

	got, err := c.GetToken(apiKey)
	if err != nil {
		t.Fatalf("Client.GetToken() | got error %v, want nil", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("Client.GetToken() | (-got +want)\n%s", diff)
	}
}

var getTokenErrors = []struct {
	name          string
	invalidURL    bool
	doError       bool
	invalidStatus bool
	badJSON       bool
}{
	{
		name:       "Invalid Request",
		invalidURL: true,
	},
	{
		name:    "Executing Request",
		doError: true,
	},
	{
		name:          "Invalid Status Code",
		invalidStatus: true,
	},
	{
		name:    "Bad JSON Response",
		badJSON: true,
	},
}

func TestClient_GetToken_Error(t *testing.T) {
	for _, tt := range getTokenErrors {
		t.Run(tt.name, func(t *testing.T) {
			apiKey := "api_key"

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = r.ParseForm()
				gotAPIKey := r.Form.Get("api_key")
				if gotAPIKey != apiKey {
					t.Fatalf("Client.GetToken() | got api token %s in request, want %s", gotAPIKey, apiKey)
				}

				if tt.invalidStatus {
					w.WriteHeader(500)
				}

				resp := []byte(getTokenResp)
				if tt.badJSON {
					resp = []byte("invalid")
				}
				_, _ = w.Write(resp)
			}))
			defer s.Close()

			c := &Client{
				baseURL: s.URL,
				hClient: http.DefaultClient,
			}

			if tt.invalidURL {
				c.baseURL = "::::"
			}

			if tt.doError {
				c.baseURL = ""
			}

			_, err := c.GetToken(apiKey)
			if err == nil {
				t.Fatalf("Client.GetToken() | got error nil, want not nil")
			}
		})
	}
}

func TestGetTokenResponse_UnmarshalJSON(t *testing.T) {
	g := &GetTokenResponse{}
	if err := g.UnmarshalJSON([]byte(getTokenResp)); err != nil {
		t.Fatalf("UnmarsharJSON() | got error %v, want nil", err)
	}

	want := &GetTokenResponse{
		Token: "1357924680",
		Expire: func() time.Time {
			tm, err := time.Parse("2006-01-02T15:04:05", "2021-09-19T01:35:57")
			if err != nil {
				t.Fatal(err)
			}
			return tm
		}(),
	}

	if diff := cmp.Diff(g, want); diff != "" {
		t.Fatalf("UnmarshalJSON() | (-got +want)\n%s", diff)
	}
}

var getTokenUnmarshalJSONError = []struct {
	name        string
	invalidJSON bool
	invalidTime bool
}{
	{
		name:        "Unmarshalling Original JSON",
		invalidJSON: true,
	},
	{
		name:        "Invalid Time",
		invalidTime: true,
	},
}

const getTokenRespInvalidTime = ` 
	{
		"token": "1357924680",
		"expire": "invalid time"
	}`

func TestGetTokenResponse_UnmarshalJSON_Error(t *testing.T) {
	for _, tt := range getTokenUnmarshalJSONError {
		t.Run(tt.name, func(t *testing.T) {
			b := []byte(getTokenResp)
			if tt.invalidJSON {
				b = []byte("invalid")
			}

			if tt.invalidTime {
				b = []byte(getTokenRespInvalidTime)
			}

			g := &GetTokenResponse{}
			if err := g.UnmarshalJSON(b); err == nil {
				t.Fatal("UnmarshalJSON() | got error nil, want not nil")
			}
		})
	}
}
