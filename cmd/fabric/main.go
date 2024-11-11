package main

import (
	"log"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/infrastructure"
)

func main() {
	c := infrastructure.NewBlockchain("./")
	c.HelloWorld()
	// c.AddPermission("doctor1", "patient1", "granted")

	permissions, err := c.QueryPermissions("doctor1", "patient1", "message")
	if err != nil {
		panic(err)
	}
	latest := entities.Permission{}
	for _, p := range permissions {
		if p.CreatedAt > latest.CreatedAt {
			latest = p
		}
	}
	log.Println("Latest permission: ", latest)
}
