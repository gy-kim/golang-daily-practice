package main

import (
	"fmt"
	"os"
	"os/user"
)

func main() {
	fmt.Println("User id:", os.Getuid())

	var u *user.User
	u, _ = user.Current()
	fmt.Print("Group ids: ")
	groupIDS, _ := u.GroupIds()
	for _, i := range groupIDS {
		fmt.Print(i, " ")
	}
	fmt.Println()
}
