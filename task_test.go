package qencode

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const createTaskResp = ` 
	{
		 "error": 0,
		 "upload_url": "https://storage.qencode.com/v1/upload_file",
		 "task_token": "471272a512d76c22665db9dcee893409"
	}`

func TestClient_CreateTask(t *testing.T) {
	token := "token"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotToken := r.Form.Get("token")
		if gotToken != token {
			t.Fatalf("Client.CreateTask() | got token %s in request, want %s", gotToken, token)
		}
		_, _ = w.Write([]byte(createTaskResp))
	}))
	defer s.Close()

	c := &Client{
		baseURL: s.URL,
		hClient: http.DefaultClient,
	}

	want := CreateTaskResponse{
		Error:     0,
		UploadURL: "https://storage.qencode.com/v1/upload_file",
		TaskToken: "471272a512d76c22665db9dcee893409",
	}

	got, err := c.CreateTask(token)
	if err != nil {
		t.Fatalf("Client.CreateTask() | got error %v, want nil", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("Client.CreateTask() | (-got +want)\n%s", diff)
	}
}

var createTaskErrors = []struct {
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

func TestClient_CreateTask_Error(t *testing.T) {
	for _, tt := range createTaskErrors {
		t.Run(tt.name, func(t *testing.T) {
			token := "token"

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			_, err := c.CreateTask(token)
			if err == nil {
				t.Fatalf("Client.CreateTask() | got error nil, want not nil")
			}
		})
	}
}

const startTaskResponse = `
	{
		 "error": 0,
		 "status_url": "https://api.qencode.com/v1/status"
	}
`

func TestClient_StartTask(t *testing.T) {
	token := "task_token"
	payload := `{"data": "data"}`
	q := StartTaskQuery{
		Query: Query{
			Source: "source",
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotToken := r.Form.Get("task_token")
		if gotToken != token {
			t.Fatalf("Client.StartTask() | got token %s in request, want %s", gotToken, token)
		}

		gotPayload := r.Form.Get("payload")
		if gotPayload != payload {
			t.Fatalf("Client.StartTask() | got payload %s in request, want %s", gotPayload, token)
		}

		d, err := json.Marshal(q)
		if err != nil {
			t.Fatal(err)
		}

		gotQuery := r.Form.Get("query")
		if gotQuery != string(d) {
			t.Fatalf("Client.StartTask() | got query %s in request, want %s", gotQuery, string(d))
		}

		_, _ = w.Write([]byte(startTaskResponse))
	}))
	defer s.Close()

	c := &Client{
		baseURL: s.URL,
		hClient: http.DefaultClient,
	}

	want := StartTaskResponse{
		Error:     0,
		StatusURL: "https://api.qencode.com/v1/status",
	}

	got, err := c.StartTask(token, payload, q)
	if err != nil {
		t.Fatalf("Client.StartTask() | got error %v, want nil", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("Client.StartTask() | (-got +want)\n%s", diff)
	}
}

var startTaskErrors = []struct {
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

func TestClient_StartTask_Error(t *testing.T) {
	for _, tt := range startTaskErrors {
		t.Run(tt.name, func(t *testing.T) {
			token := "task_token"
			payload := `{"data": "data"}`
			q := StartTaskQuery{
				Query: Query{
					Source: "source",
				},
			}

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			_, err := c.StartTask(token, payload, q)
			if err == nil {
				t.Fatalf("Client.StartTask() | got error nil, want not nil")
			}
		})
	}
}
