package main

import (
	api "emurobot/pkg/api"
	"fmt"
	"os"
)

func main() {
	go api.InitServerWithWeb()

	fmt.Println(os.Hostname())

	for {

	}
}
