package main

import (
	"fmt"
	"gogetent/os/user"
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

var mu = sync.Mutex{}

func iterateUsers(listNum int) {
	mu.Lock()
	l := make([]string, 0, 10)
	_ = user.IterateUsers(func(u *user.User) error {
		l = append(l, u.Username)
		return nil
	})
	fmt.Printf("Goroutine %d user names: %+v \n\n", listNum, l)
	wg.Done()
	mu.Unlock()
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
