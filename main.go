package main

import (
	"fmt"

	"github.com/s4bb4t/slib/handle"
)

func main() {
	handle.ChangeConfig(1, 3, true)
	fmt.Println(handle.Attack("POST", "https://easydev.club/api/v1/todos", []byte(`{"title":"lol"}`)))
}
