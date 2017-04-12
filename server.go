package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Entrada struct {
	Sitio      string
	User       string
	Password   string
	Comentario string
}

type Usuario struct {
	Sal       string
	MasterKey string
	Username  string
	Lista     []Entrada
}

func login(user string, password string) bool {
	usuarios := map[int]Usuario{}
	file, err := os.Open("bd.json")
	chk(err)
	defer file.Close()
	str, err := ioutil.ReadAll(file)
	if erru := json.Unmarshal(str, &usuarios); erru != nil {
		panic(erru)
	}
	for key, value := range usuarios {
		if value.Username == user {
			sha_512 := sha512.New()
			sha_512.Write([]byte(password))
			pass2 := encode64(sha_512.Sum([]byte(value.Sal)))
			if value.MasterKey == pass2 {
				fmt.Print(key)
				return true
			}
		}
	}
	return false
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

func registro(user string, password string) bool {
	usuarios := map[int]Usuario{}
	file, err := os.Open("bd.json")
	chk(err)
	defer file.Close()
	str, err := ioutil.ReadAll(file)
	if erru := json.Unmarshal(str, &usuarios); erru != nil {
		panic(erru)
	}
	//comprobamos si existe el usuario
	for key, value := range usuarios {
		if value.Username == user {
			fmt.Println(key)
			return false
		}
	}
	salG, error := GenerateRandomString(10)
	chk(error)

	sha_512 := sha512.New()
	sha_512.Write([]byte(password))
	pass2 := encode64(sha_512.Sum([]byte(salG)))

	var lista []Entrada
	usuario_nuevo := Usuario{Sal: salG, MasterKey: pass2, Username: user, Lista: lista}
	usuarios[len(usuarios)+1] = usuario_nuevo //añadimos el nuevo usuario al map
	usuarios_json, err := json.Marshal(usuarios)
	if err != nil {
		fmt.Println("Error marshal: ", err)
	}

	ioutil.WriteFile("bd.json", usuarios_json, 0644)

	return true
}

func chk(e error) {
	if e != nil {
		panic(e)
	}
}

type resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
}

func response(w io.Writer, ok bool, msg string) {
	r := resp{Ok: ok, Msg: msg}    // formateamos respuesta
	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	// Redirect the incoming HTTP request. Note that "127.0.0.1:8081" will only work if you are accessing the server from your local machine.
	http.Redirect(w, r, "https://127.0.0.1:8081"+r.RequestURI, http.StatusMovedPermanently)
}

func handler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()                              // es necesario parsear el formulario
	w.Header().Set("Content-Type", "text/plain") // cabecera estándar

	switch req.Form.Get("cmd") { // comprobamos comando desde el cliente
	case "Registro":
		{
			mensaje := ""
			if registro(req.Form.Get("Usuario"), req.Form.Get("Password")) {

				fmt.Println("Registro ok")
				mensaje = "Usuario: " + req.Form.Get("Usuario") + ", Password: " + req.Form.Get("Password")
				response(w, true, mensaje)
			} else {
				fmt.Println("Error en el registro")
				mensaje = "Usuario ya existe"
				response(w, false, mensaje)
			}
		}

	case "Login":
		{
			mensaje := ""
			if login(req.Form.Get("Usuario"), req.Form.Get("Password")) {
				fmt.Println("Login ok")
				mensaje = "Usuario: " + req.Form.Get("Usuario") + ", Password: " + req.Form.Get("Password")
				response(w, true, mensaje)
			} else {
				fmt.Println("Login error")
				mensaje = "Usuario y/o password incorrecto"
				response(w, false, mensaje)
			}

		}
	default:
		response(w, false, "Comando inválido")
	}
}

func main() {

	fmt.Println("Servidor emcendido en el puerto 10443...")
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handler))

	srv := &http.Server{Addr: ":10443", Handler: mux}

	go func() {
		if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // espera señal SIGINT
	log.Println("Apagando servidor ...")

	// apagar servidor de forma segura
	//ctx, fnc := context.WithTimeout(context.Background(), 5*time.Second)
	//fnc()
	//srv.Shutdown(ctx)

	log.Println("Servidor detenido correctamente")
}
