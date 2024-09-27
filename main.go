package main

import (
	"fmt"

	"github.com/sabbatD/slib/handle"
)

func main() {
	handle.ChangeConfig(1, 3, 1)
	fmt.Println(handle.Attack("GET", "https://easydev.club/api/v1/todos"))
}
