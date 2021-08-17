package main

import (
	"fmt"
	"log"
	"os/user"
	"reflect"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}
var mu = sync.Mutex{}

func main() {
	fmt.Println("Started users groups iteration process")

	users := iterateUsers()
	groups := iterateGroups()

	// Print out iterated users/groups
	userNames := make([]string, 0, len(users))
	for _, u := range users {
		userNames = append(userNames, u.Username)
	}
	groupNames := make([]string, 0, len(groups))
	for _, group := range groups {
		groupNames = append(groupNames, group.Name)
	}

	fmt.Printf("Usernames: %+v \n\nGroupnames: %+v \n\n", strings.Join(userNames, ","), strings.Join(groupNames, ","))

	num := 5
	fmt.Printf("Running %d goroutines for users and groups \n", num)

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			u := iterateUsers()
			if reflect.DeepEqual(users, u) {
				fmt.Printf("Goroutine [Users#%d] iterated all %d users \n", i, len(u))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			g := iterateGroups()
			if reflect.DeepEqual(groups, g) {
				fmt.Printf("Goroutine [Groups#%d] iterated all %d groups \n", i, len(g))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

}

func iterateUsers() []*user.User {
	mu.Lock()
	l := make([]*user.User, 0, 10)
	err := user.IterateUsers(func(u *user.User) error {
		l = append(l, u)
		return nil
	})
	if err != nil {
		log.Fatalf("error occurred while iterating users: %v \n", err)
	}
	mu.Unlock()
	return l
}

func iterateGroups() []*user.Group {
	mu.Lock()
	l := make([]*user.Group, 0, 10)
	err := user.IterateGroups(func(u *user.Group) error {
		l = append(l, u)
		return nil
	})
	if err != nil {
		log.Fatalf("error occurred while iterating groups: %v \n", err)
	}
	mu.Unlock()
	return l
}
