package main

import (
	"gogetent/os/user"

	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	num := 3

	for i := 0; i < num; i++ {
		go iterateUsers(i)
		wg.Add(1)
	}
	wg.Wait()

	fmt.Println("Starting to iterate over groups")
	for i := 0; i < num; i++ {
		go iterateGroups(i)
		wg.Add(1)
	}
	wg.Wait()
}

func iterateUsers(listNum int) {
	l := make([]string, 0, 10)
	_ = user.IterateUsers(func(u *user.User) error {
		l = append(l, u.Username)
		return nil
	})
	fmt.Printf("Goroutine %d user names: %+v \n\n", listNum, l)
	wg.Done()
}

func iterateGroups(listNum int) {
	l := make([]string, 0, 10)
	_ = user.IterateGroups(func(u *user.Group) error {
		l = append(l, u.Name)
		return nil
	})
	fmt.Printf("Goroutine %d group names: %+v \n\n", listNum, l)
	wg.Done()
}
