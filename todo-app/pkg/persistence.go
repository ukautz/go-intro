package todo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Persistence is a storage for todos
type Persistence interface {

	// Create stores a new Todo and returns the ID
	Create(todo Todo) (string, error)

	// Delete removes a single Todo identified by it's ID. Returns os.ErrNotExist if not found
	Delete(id string) error

	// Get fetches a single Todo identified by it's ID. Returns os.ErrNotExist if not found
	Get(id string) (*Todo, error)

	// List returns all Todos
	List() ([]Todo, error)
}

// DirectoryPersistence implements Persistence with a local file system directory
type DirectoryPersistence string

// Create stores Todo in <directory>/<id>.json file
func (p DirectoryPersistence) Create(todo Todo) (string, error) {
	if todo.ID == "" {
		todo.ID = uuid.New().String()
		todo.Created = time.Now()
	}
	encoded, err := json.Marshal(todo)
	if err != nil {
		return "", err
	}

	path := p.path(todo.ID)
	err = ioutil.WriteFile(path, encoded, 0640)
	if err != nil {
		return "", err
	}

	return todo.ID, nil
}

// Delete removes <directory>/<id>.json file
func (p DirectoryPersistence) Delete(id string) error {
	return os.Remove(p.path(id))
}

// Get reads Todo from <directory>/<id>.json file
func (p DirectoryPersistence) Get(id string) (*Todo, error) {
	return p.read(p.path(id))
}

// List reads all Todos from <id>.json files in <directory>
func (p DirectoryPersistence) List() ([]Todo, error) {
	todos := make([]Todo, 0)
	err := filepath.Walk(string(p), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil
		} else if filepath.Ext(path) != ".json" {
			return nil
		}

		todo, err := p.read(path)
		if err != nil {
			return err
		}

		todos = append(todos, *todo)
		return nil
	})

	return todos, err
}

func (p DirectoryPersistence) path(id string) string {
	return filepath.Join(string(p), id+".json")
}

func (p DirectoryPersistence) read(path string) (*Todo, error) {
	encoded, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var todo Todo
	if err = json.Unmarshal(encoded, &todo); err != nil {
		return nil, err
	}

	return &todo, nil
}
