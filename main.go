package main

import (
	"./os/user"
	"fmt"
)

func main(){
	us, err := user.GetAllUsers()
	fmt.Println(len(us))
	if err != nil{
		panic(err)
	}
	for _, u := range us{
		fmt.Printf("%+v \n", u)
	}
	gs, err := user.GetAllGroups()
	if err != nil{
		panic(err)
	}
	for _, g := range gs{
		fmt.Printf("%+v \n", g)
	}
}
