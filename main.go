package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
)

func main() {
	args := os.Args

	if len(args) == 2 {
		db := args[1]

		switch db {
		case "passwd":
			err := user.IterateUsers(func(u *user.User) error {
				fmt.Printf("%s:x:%s:%s:%s:%s\n", u.Username, u.Uid, u.Gid, u.Name, u.HomeDir)
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		case "group":
			err := user.IterateGroups(func(g *user.Group) error {
				fmt.Printf("%s:x:%s:\n", g.Name, g.Gid)
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		os.Exit(0)
	}

	log.Fatalf("Please provide passwd or group as second argument!")
}
