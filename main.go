package main

import (
	"fmt"

	"github.com/sabbatD/slib/handle"
)

func main() {
	handle.ChangeConfig(1, 3, true)
	fmt.Println(handle.Attack("POST", "https://easydev.club/api/v1/todos", []byte(`{"title":"lol"}`)))
}
