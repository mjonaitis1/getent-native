package main

import (
	"./os/user"
	"fmt"
)

func main(){
	u, _ := user.Lookup("mjonaitis")
	fmt.Printf("%+v", u)
}
