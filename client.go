package main

import (
	"crypto/sha512"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Entrada struct {
	Sitio      string
	User       string
	Password   string
	Comentario string
}

type respGeneral struct {
	Ok  bool
	Msg string
}

type respSesion struct {
	Ok     bool
	Msg    string
	Sesion bool
}

type respEntrada struct {
	Ok           bool
	ValorEntrada Entrada
}

func chk(e error) {
	if e != nil {
		panic(e)
	}
}

var usuario string

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
func sendPost(data url.Values) []byte {
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

	return body
}

func login() bool {
	var usuarioL string
	fmt.Println("--- Iniciar sesión: ---")
	fmt.Print("Introduce usuario: ")
	fmt.Scanf("%s\n", &usuarioL)
	fmt.Print("Introduce password: ")
	password, err := gopass.GetPasswd()
	chk(err)
	data := url.Values{} // estructura para contener los valores

	sha_512 := sha512.New()
	sha_512.Write([]byte(password))
	pass2 := encode64(sha_512.Sum(nil))

	data.Set("cmd", "Login")
	data.Set("Usuario", usuarioL) // comando (string)
	data.Set("Password", pass2)   // usuario (string)

	body := sendPost(data)
	respuesta := respGeneral{}

	if erru := json.Unmarshal(body, &respuesta); erru != nil {
		panic(erru)
	}
	fmt.Println(respuesta.Msg)
	if respuesta.Ok {
		usuario = usuarioL
		return true
	}
	return false
}

func logout() bool {
	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "Logout")
	data.Set("Usuario", usuario)

	body := sendPost(data)
	respuesta := respGeneral{}

	if erru := json.Unmarshal(body, &respuesta); erru != nil {
		panic(erru)
	}

	fmt.Println(respuesta.Msg)
	return respuesta.Ok
}
func comprobarSesion(usuario string) (resultado bool, mensaje string) {
	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "Session")
	data.Set("Usuario", usuario)
	body := sendPost(data)
	respuesta := respSesion{}

	if erru := json.Unmarshal(body, &respuesta); erru != nil {
		panic(erru)
	}
	return respuesta.Ok, respuesta.Msg
}

func pedirPassword() string {
	var opcion int
	var pass2 string
	for !(opcion < 3 && opcion > 0) {
		fmt.Println("1-Introducir contraseña")
		fmt.Println("2-Generar contraseña aleatoria")

		fmt.Scanf("%d\n", &opcion)

		switch opcion {
		case 1:
			{
				fmt.Print("Introduce password: ")
				password, err := gopass.GetPasswd()
				chk(err)
				pass2 = encode64(password)
			}
		case 2:
			{
				var tam int
				fmt.Println("Introduzca el tamaño de la contraseña generada")
				fmt.Scanf("%d\n", &tam)
				var contr = RandomPass(tam)
				fmt.Print("La contraseña es: ")
				fmt.Println(string(contr))
				pass2 = encode64(contr)
			}
		}
	}
	return pass2
}

func add() bool {

	fmt.Println("--- Añadir nueva cuenta ---")
	fmt.Print("Introduce sitio web: ")
	var sitio string
	fmt.Scanf("%s\n", &sitio)
	fmt.Print("Introduce usuario: ")
	var usuariositio string
	fmt.Scanf("%s\n", &usuariositio)

	var pass2 = pedirPassword()
	fmt.Print("Introduce un comentario: ")
	var comentario string
	fmt.Scanf("%s\n", &comentario)

	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "Add")
	data.Set("Usuario", usuario)
	data.Set("Sitio", sitio)
	data.Set("Usuariositio", usuariositio)
	data.Set("Password", pass2)
	data.Set("Comentario", comentario)
	body := sendPost(data)
	respuesta := respGeneral{}
	if erru := json.Unmarshal(body, &respuesta); erru != nil {
		panic(erru)
	}

	fmt.Println(respuesta.Msg)

	return respuesta.Ok

}

func view() bool {

	fmt.Println("--- Ver una cuenta ---")
	fmt.Print("Introduce sitio web: ")
	var sitio string
	fmt.Scanf("%s\n", &sitio)
	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "View")
	data.Set("Usuario", usuario)
	data.Set("Sitio", sitio)

	body := sendPost(data)

	respuesta := respEntrada{}

	if erru := json.Unmarshal(body, &respuesta); erru != nil {
		panic(erru)
	}
	entrada := respuesta.ValorEntrada
	if entrada.User != "." {
		fmt.Println(respuesta.Ok)
		fmt.Println("Sitio: ")
		fmt.Println(sitio)
		fmt.Println("User: ")
		fmt.Println(entrada.User)
		fmt.Println("Password: ")
		fmt.Println(entrada.Password)
		fmt.Println("Comentario: ")
		fmt.Println(entrada.Comentario)
		return true
	} else {
		fmt.Println("Entrada no añadida aun")
	}
	return false

}

func delete() bool {

	fmt.Print("Introduce sitio web: ")
	var sitio string
	fmt.Scanf("%s\n", &sitio)
	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "Delete")
	data.Set("Usuario", usuario)
	data.Set("Sitio", sitio)

	body := sendPost(data)

	respuesta := respGeneral{}

	if erru := json.Unmarshal(body, &respuesta); erru != nil {
		panic(erru)
	}
	fmt.Println(respuesta.Msg)
	return respuesta.Ok

}
func edit() bool {

	fmt.Println("--- Editar una cuenta ---")
	fmt.Print("Introduce sitio web: ")
	var sitio string
	fmt.Scanf("%s\n", &sitio)
	data := url.Values{} // estructura para contener los valores
	data.Set("cmd", "Edit?")
	data.Set("Usuario", usuario)
	data.Set("Sitio", sitio)

	body := sendPost(data)

	respuesta := respGeneral{}

	if erru := json.Unmarshal(body, &respuesta); erru != nil {
		panic(erru)
	}
	if respuesta.Ok {
		fmt.Print("Introduce usuario: ")
		var usuariositio string
		fmt.Scanf("%s\n", &usuariositio)
		var pass2 = pedirPassword()

		fmt.Print("Introduce un comentario: ")
		var comentario string
		fmt.Scanf("%s\n", &comentario)

		data := url.Values{} // estructura para contener los valores
		data.Set("cmd", "Edit")
		data.Set("Usuario", usuario)
		data.Set("Sitio", sitio)
		data.Set("Usuariositio", usuariositio)
		data.Set("Password", pass2)
		data.Set("Comentario", comentario)
		body := sendPost(data)
		respuesta := respGeneral{}
		if erru := json.Unmarshal(body, &respuesta); erru != nil {
			panic(erru)
		}
		if respuesta.Ok {
			fmt.Println("Entrada editada")
			return true
		}
	} else {
		fmt.Println("Entrada no añadida aun")
	}
	return false
}
func menuprincipal() {

	var opcion int
	for opcion != 5 {
		fmt.Println("--- Sesión iniciada ---")
		fmt.Println("Elige la opción que desea realizar: ")
		fmt.Println("1 - Consultar una cuenta.")
		fmt.Println("2 - Añadir nueva cuenta.")
		fmt.Println("3 - Eliminar una cuenta.")
		fmt.Println("4 - Editar una cuenta.")
		fmt.Println("5 - Salir(Logout).")
		fmt.Print("Opción: ")
		fmt.Scanf("%d\n", &opcion)
		comp, mens := comprobarSesion(usuario)
		if !comp {
			fmt.Println(mens)
			opcion = 5
		}
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
				delete()
			}
		case 4:
			{
				edit()
			}
		case 5:
			{
				logoutok := logout()
				if logoutok {
					usuario = ""
				}
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

func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	chk(err)                                     // comprobamos el error
	return b                                     // devolvemos los datos originales
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
	body := sendPost(data)
	respuesta := respGeneral{}
	if erru := json.Unmarshal(body, &respuesta); erru != nil {
		panic(erru)
	}
	fmt.Println(respuesta.Msg)
	return respuesta.Ok

}

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomPass(strlen int) []byte {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return result
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
