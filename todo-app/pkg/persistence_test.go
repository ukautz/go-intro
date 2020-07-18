package todo_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	todo "github.com/ukautz/go-intro/todo-app/pkg"
)

func TestDirectoryPersistence_Create(t *testing.T) {
	p := createTestDirectoryPersistence(t)
	id, err := p.Create(todo.Todo{
		Title:       "the-title",
		Description: "the-description",
		UserID:      "u01",
	})
	require.NoError(t, err)
	require.NotEmpty(t, id)

	testFile := filepath.Join(testPersistenceDir, fmt.Sprintf("%s.json", id))
	defer os.Remove(testFile)

	raw, err := ioutil.ReadFile(testFile)
	require.NoError(t, err)

	var td todo.Todo
	require.NoError(t, json.Unmarshal(raw, &td))

	assert.Equal(t, id, td.ID)
	assert.Equal(t, "u01", td.UserID)
	assert.Equal(t, "the-title", td.Title)
	assert.Equal(t, "the-description", td.Description)
}

func TestDirectoryPersistence_Delete(t *testing.T) {
	storePath := assertJSONTodoFile(t, 2)
	defer os.Remove(storePath)

	p := createTestDirectoryPersistence(t)
	err := p.Delete("todo-02")
	require.NoError(t, err)

	_, err = os.Stat(storePath)
	assert.True(t, os.IsNotExist(err))

	err = p.Delete("todo-02")
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestDirectoryPersistence_Get(t *testing.T) {
	defer os.Remove(assertJSONTodoFile(t, 2))

	p := createTestDirectoryPersistence(t)
	td, err := p.Get("todo-02")
	require.NoError(t, err)
	assert.Equal(t, &todo.Todo{
		ID:          "todo-02",
		Title:       "todo 02",
		Description: "the todo number 02",
		UserID:      "u02",
		Created:     time.Date(2010, 11, 12, 13, 14, 15, 0, time.UTC),
	}, td)
}

func TestDirectoryPersistence_List(t *testing.T) {
	defer os.Remove(assertJSONTodoFile(t, 1))
	defer os.Remove(assertJSONTodoFile(t, 3))
	defer os.Remove(assertJSONTodoFile(t, 5))

	p := createTestDirectoryPersistence(t)

	tds, err := p.List()
	require.NoError(t, err)
	require.Len(t, tds, 3)
	assert.Equal(t, []todo.Todo{
		{
			ID:          "todo-01",
			Title:       "todo 01",
			Description: "the todo number 01",
			UserID:      "u01",
			Created:     time.Date(2010, 11, 12, 13, 14, 15, 0, time.UTC),
		},
		{
			ID:          "todo-03",
			Title:       "todo 03",
			Description: "the todo number 03",
			UserID:      "u03",
			Created:     time.Date(2010, 11, 12, 13, 14, 15, 0, time.UTC),
		},
		{
			ID:          "todo-05",
			Title:       "todo 05",
			Description: "the todo number 05",
			UserID:      "u05",
			Created:     time.Date(2010, 11, 12, 13, 14, 15, 0, time.UTC),
		},
	}, tds)
}

var (
	testPersistenceDir = filepath.Join("fixtures", "store")
)

// createTestDirectoryPersistence assures the testing directory is existing and
// returns a new DirectoryPersistence instance of it
func createTestDirectoryPersistence(t *testing.T) todo.DirectoryPersistence {
	if err := os.MkdirAll(testPersistenceDir, 0755); err != nil {
		t.Fatal(err)
	}
	return todo.DirectoryPersistence(testPersistenceDir)
}

// assertJSONTodoFile creates a test fixtures JSON file and returns it's absolute path
func assertJSONTodoFile(t *testing.T, num int) string {

	id := fmt.Sprintf("%02d", num)
	fileName := "todo-" + id + ".json"
	storePath := filepath.Join(testPersistenceDir, fileName)

	encoded := `{"id":"todo-:num:","title":"todo :num:","description":"the todo number :num:","created":"2010-11-12T13:14:15Z","user_id":"u:num:"}`
	encoded = strings.ReplaceAll(encoded, ":num:", id)

	err := ioutil.WriteFile(storePath, []byte(encoded), 0644)
	require.NoError(t, err)
	return storePath
}
