package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func menu() int {
	fmt.Println("Cliente de la aplicación")
	fmt.Println("Elige la opción que desea realizar: ")
	fmt.Println("1 - Iniciar sesión.")
	fmt.Println("2 - Registrarse.")
	fmt.Println("3 - Salir.")
	var opcion int
	for opcion < 1 || opcion > 3 {
		fmt.Print("Opción: ")
		fmt.Scanf("%d\n", &opcion)
		if opcion < 1 || opcion > 3 {
			fmt.Println("Opción incorrecta. Debe ser un valor entre 1 y 3")
		}
	}

	return opcion

}

func login() bool {
	fmt.Println("--- Iniciar sesión: ---")
	fmt.Print("Introduce usuario: ")
	var usuario string
	fmt.Scanf("%s\n", &usuario)
	fmt.Print("Introduce password: ")
	var password string
	fmt.Scanf("%s\n", &password)
	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "Login")
	data.Set("Usuario", usuario)   // comando (string)
	data.Set("Password", password) // usuario (string)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	r, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	chk(err)
	io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	fmt.Println()
	return true

}

func registro() bool {
	fmt.Println("--- Registrarse: ---")
	fmt.Print("Introduce usuario: ")
	var usuario string
	fmt.Scanf("%s\n", &usuario)
	fmt.Print("Introduce password: ")
	var password string
	fmt.Scanf("%s\n", &password)
	data := url.Values{} // estructura para contener los valores

	data.Set("cmd", "Registro")
	data.Set("Usuario", usuario)   // comando (string)
	data.Set("Password", password) // usuario (string)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	r, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	chk(err)
	io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	return true
}

func main() {

	opcion := menu()

	switch opcion {
	case 1:
		{
			login()
		}
	case 2:
		{
			registro()
		}
	default:
		fmt.Println("Cerrando cliente...")
	}

}
