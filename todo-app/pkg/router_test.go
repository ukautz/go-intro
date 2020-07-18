package todo_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	todo "github.com/ukautz/go-intro/todo-app/pkg"
)

func TestRouter_ServeHTTP_RejectInvalidCredentials(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/todo", nil)
	req.SetBasicAuth("invalid", "invalid")

	router := testNewRouter()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	require.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestRouter_ServeHTTP_Create(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/todo", bytes.NewBuffer([]byte(`{"title":"the-title", "description": "the-description"}`)))
	req.SetBasicAuth("the-user", "the-pass")

	router := testNewRouter()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)

	ret, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	out := make(map[string]string)
	require.NoError(t, json.Unmarshal(ret, &out))
	require.Contains(t, out, "id")
	assert.NotEmpty(t, out["id"])
}

func TestRouter_ServeHTTP_List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/todo", nil)
	req.SetBasicAuth("the-user", "the-pass")

	router := testNewRouter()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)

	ret, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	out := make([]todo.Todo, 0)
	require.NoError(t, json.Unmarshal(ret, &out))
	assert.Equal(t, []todo.Todo{
		{
			ID:    "todo-01",
			Title: "todo 01",
		},
		{
			ID:    "todo-02",
			Title: "todo 02",
		},
	}, out)
}

func TestRouter_ServeHTTP_Delete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/todo/todo-01", nil)
	req.SetBasicAuth("the-user", "the-pass")

	router := testNewRouter()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)

	ret, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, `{"id":"todo-01"}`, strings.TrimSpace(string(ret)))
}

func TestRouter_ServeHTTP_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/todo/todo-01", nil)
	req.SetBasicAuth("the-user", "the-pass")

	router := testNewRouter()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)

	ret, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	out := todo.Todo{}
	require.NoError(t, json.Unmarshal(ret, &out))
	assert.Equal(t, todo.Todo{
		ID:    "todo-01",
		Title: "todo 01",
	}, out)
}

func TestRouter_ServeHTTP_Fallback(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("the-user", "the-pass")

	router := testNewRouter()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func testNewRouter() todo.Router {
	return todo.Router{
		Authentication: testAuthentication{"the-user": "the-pass"},
		Persistence: testPersistence{
			"todo-01": {
				ID:    "todo-01",
				Title: "todo 01",
			},
			"todo-02": {
				ID:    "todo-02",
				Title: "todo 02",
			},
		},
	}
}

var (
	// testIDCount contains an incrementing counter to generate unique test IDS
	testIDCount = 0
)

// testID creates an incrementing ID for testing
func testID() string {
	testIDCount++
	return fmt.Sprintf("u%03d", testIDCount)
}

type testPersistence map[string]todo.Todo

func (p testPersistence) Create(td todo.Todo) (string, error) {
	td.ID = testID()
	p[td.ID] = td
	return td.ID, nil
}

func (p testPersistence) Delete(id string) error {
	if _, ok := p[id]; ok {
		delete(p, id)
		return nil
	}
	return errors.New("not found")
}

func (p testPersistence) Get(id string) (*todo.Todo, error) {
	if td, ok := p[id]; ok {
		return &td, nil
	}
	return nil, errors.New("not found")
}

func (p testPersistence) List() ([]todo.Todo, error) {
	ids := make([]string, 0)
	for id := range p {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	tds := make([]todo.Todo, len(ids))
	for i, id := range ids {
		tds[i] = p[id]
	}
	return tds, nil
}

type testAuthentication map[string]string

func (a testAuthentication) Authenticate(req *http.Request) (userID string, err error) {
	user, pass, hasBasic := req.BasicAuth()
	if !hasBasic {
		return "", errors.New("missing basic credentials")
	}
	if known, has := a[user]; has && pass == known {
		return user, nil
	}
	return "", todo.NotAllowedError
}
