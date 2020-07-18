package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	todo "github.com/ukautz/go-intro/todo-app/pkg"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "server"
	app.Usage = "HTTP API for todos"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "storage-directory",
			Aliases: []string{"d"},
			Usage:   "Path to directory to store todos",
			Value:   filepath.Join("data", "store"),
		},
		&cli.StringFlag{
			Name:    "users",
			Aliases: []string{"u"},
			Usage:   "Path to JSON file containing user credentials",
			Value:   filepath.Join("data", "users.json"),
		},
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a"},
			Usage:   "Listen address for the HTTP server",
			Value:   "127.0.0.1:12345",
		},
		&cli.StringFlag{
			Name:    "path-prefix",
			Aliases: []string{"p"},
			Usage:   "Prefix for ",
			Value:   "/v1",
		},
	}

	app.Action = func(c *cli.Context) error {
		listenAddr := c.String("address")
		routePrefix := c.String("path-prefix")

		// init storage
		store := todo.DirectoryPersistence(c.String("storage-directory"))

		// load users for authentication
		usersFile := c.String("users")
		auth, err := todo.LoadAuthenticationFromJSON(usersFile)
		if err != nil {
			return err
		}

		// setup router
		router := todo.Router{
			Prefix:         routePrefix,
			Authentication: auth,
			Persistence:    store,
		}

		// run server
		log.Printf("Starting API server at http://%s%s, storage directory: %s",
			listenAddr, routePrefix, store)
		return http.ListenAndServe(listenAddr, router)
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
