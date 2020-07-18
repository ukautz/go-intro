package todo_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	todo "github.com/ukautz/go-intro/todo-app/pkg"
)

func TestUsersAuthentication_Authenticate(t *testing.T) {
	auth := todo.UsersAuthentication{
		{ID: "u01", Name: "alice", Password: "secret1"},
		{ID: "u02", Name: "bob", Password: "secret2"},
	}

	expects := []struct {
		name    string
		request *http.Request
		id      string
		allowed bool
	}{
		{"missing basic auth forbidden", createBasicAuthTestRequest("", ""), "", false},
		{"unknown credentials forbidden", createBasicAuthTestRequest("foo", "bar"), "", false},
		{"invalid credentials forbidden", createBasicAuthTestRequest("alice", "invalid"), "", false},
		{"allow valid user u01", createBasicAuthTestRequest("alice", "secret1"), "u01", true},
		{"allow valid user u02", createBasicAuthTestRequest("bob", "secret2"), "u02", true},
	}

	for _, expect := range expects {
		//expect := expect
		t.Run(expect.name, func(t *testing.T) {
			//t.Parallel()
			userID, err := auth.Authenticate(expect.request)
			if expect.allowed {
				assert.NoError(t, err)
				assert.Equal(t, expect.id, userID)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func createBasicAuthTestRequest(user, pass string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:12345/bla", nil)
	if user != "" {
		req.SetBasicAuth(user, pass)
	}
	return req
}
