package main

import (
	"fmt"
	"ioutil"
	"net/http"
)

func main() {
	fmt.Println("Cliente:")
	resp, err := http.Get("http:localhost:8080/monkeys")
	if err != nil {
		fmt.Println("Error en petición")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
}
