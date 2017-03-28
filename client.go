package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
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
	resp, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body) // mostramos el cuerpo de la respuesta (es un reader)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Cuerpo", body)
	return true

}

func menuprincipal() {

	var opcion int
	for opcion != 4 {
		fmt.Println("--- Sesión iniciada ---")
		fmt.Println("Elige la opción que desea realizar: ")
		fmt.Println("1 - Consultar una cuenta.")
		fmt.Println("2 - Añadir nueva cuenta.")
		fmt.Println("3 - Eliminar una cuenta.")
		fmt.Println("4 - Salir.")
		fmt.Print("Opción: ")
		fmt.Scanf("%d\n", &opcion)

		switch opcion {
		case 1:
			{
				//Consultar
			}
		case 2:
			{
				//Añadir
			}
		case 3:
			{
				//Eliminar
			}
		case 4:
			{
				//Salir. No hace nada
			}
		default:
			{
				fmt.Println("Opción incorrecta. Debe ser un valor entre 1 y 3")
			}
		}
	}
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

	salir := false
	for salir == false {
		opcion := menu()

		switch opcion {
		case 1:
			{
				if logueado := login(); logueado {
					menuprincipal()
				}
			}
		case 2:
			{
				registro()
			}
		default:
			salir = true
			fmt.Println("Cerrando cliente...")
		}
	}

}
