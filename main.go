package main

import (
	"fmt"

	"github.com/sabbatD/slib/handle"
)

func main() {
	handle.ChangeConfig(400, 3, false)
	fmt.Println(handle.Attack("GET", "https://easydev.club/api/v1/todos"))
}
