package main

import (
	"crypto/sha512"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"io/ioutil"
	"github.com/howeyc/gopass"
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
func sendPost(data url.Values)string{
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	chk(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		fmt.Println("errorination happened reading the body", err)
		panic("Error respuesta Post")
	}
	a := string(body[:])
	return a
}

func login() bool {
	fmt.Println("--- Iniciar sesión: ---")
	fmt.Print("Introduce usuario: ")
	var usuario string
	fmt.Scanf("%s\n", &usuario)
	fmt.Print("Introduce password: ")
	password, err := gopass.GetPasswd()
	chk(err)
	data := url.Values{} // estructura para contener los valores

	sha_512 := sha512.New()
	sha_512.Write([]byte(password))
	pass2 := encode64(sha_512.Sum(nil))

	data.Set("cmd", "Login")
	data.Set("Usuario", usuario) // comando (string)
	data.Set("Password", pass2)  // usuario (string)

	var respuesta string
	respuesta = sendPost(data)
	fmt.Println(respuesta)
	fmt.Println(respuesta[0:7])
	/*if respuesta[6] != "t" {
		println("\nnError de identificación\n")
		return false
	} else {*/
		return true
	
}

func logout() {
	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "Logout")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	chk(err)
	defer resp.Body.Close()
}

func add() bool {
	fmt.Println("--- Añadir nueva cuenta ---")
	fmt.Print("Introduce sitio web: ")
	var sitio string
	fmt.Scanf("%s\n", &sitio)
	fmt.Print("Introduce usuario: ")
	var usuario string
	fmt.Scanf("%s\n", &usuario)
	fmt.Print("Introduce password: ")
	var password string
	fmt.Scanf("%s\n", &password)
	fmt.Print("Introduce un comentario: ")
	var comentario string
	fmt.Scanf("%s\n", &comentario)

	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "Add")
	data.Set("Sitio", sitio)
	data.Set("Usuario", usuario)
	data.Set("Password", password)
	data.Set("Comentario", comentario)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	chk(err)
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return true

}

func view() bool {
	fmt.Println("--- Ver una cuenta ---")
	fmt.Print("Introduce sitio web: ")
	var sitio string
	fmt.Scanf("%s\n", &sitio)
	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "View")
	data.Set("Sitio", sitio)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
	chk(err)
	defer resp.Body.Close()
	println("\n")
	io.Copy(os.Stdout, resp.Body)
	println("\n")
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
		fmt.Println("4 - Salir(Logout).")
		fmt.Print("Opción: ")
		fmt.Scanf("%d\n", &opcion)

		switch opcion {
		case 1:
			{
				view()
			}
		case 2:
			{
				add()
			}
		case 3:
			{
				//Eliminar
			}
		case 4:
			{
				logout()
			}
		default:
			{
				fmt.Println("Opción incorrecta. Debe ser un valor entre 1 y 3")
			}
		}
	}
}

func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

func registro() bool {
	fmt.Println("--- Registrarse: ---")
	fmt.Print("Introduce usuario: ")
	var usuario string
	fmt.Scanf("%s\n", &usuario)
	fmt.Print("Introduce password: ")
	password, err := gopass.GetPasswd()
	chk(err)

	sha_512 := sha512.New()
	sha_512.Write([]byte(password))
	pass2 := encode64(sha_512.Sum(nil))

	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "Registro")
	data.Set("Usuario", usuario) // comando (string)
	data.Set("Password", pass2)  // usuario (string)
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
