package main

import (
	"fmt"

	"github.com/sabbatD/slib/handle"
)

func main() {
	handle.ChangeConfig(100, 3, 1)
	fmt.Println(handle.Attack("GET", "https://easydev.club/api/v1/todos"))
}
