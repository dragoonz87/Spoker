package main

import (
	"fmt"
	"spoker/server"
)

func main() {
	fmt.Println("This is Spoker!")

    server.Listen(8080);
}

