package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"crypto/tls"
)

func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	fmt.Println("Cliente de la aplicación")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	fmt.Print("Introduce usuario: ")
	var usuario string
	fmt.Scanf("%s", &usuario)
	fmt.Print("Introduce password: ")
	var password string
	fmt.Scanf("%s", &password)
	// ** ejemplo de registro
	data := url.Values{}             // estructura para contener los valores

	data.Set("cmd", "Login")
	data.Set("Usuario", usuario)          // comando (string)
	data.Set("Password", password) // usuario (string)

	r, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	chk(err)
	io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	fmt.Println()
}
