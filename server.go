package main

import (
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
	file, err := os.Open("bd.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	str, err := ioutil.ReadAll(file)
	if erru := json.Unmarshal(str, &usuarios); erru != nil {
		panic(erru)
	}
	for key, value := range usuarios {
		if value.Username == user && value.MasterKey == password {
			return true
		}
	}
	return false
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
	case "hola": // ** registro
		response(w, true, "Hola "+req.Form.Get("mensaje"))
	case "Login":
		{
			mensaje := ""
			if login(req.Form.Get("Usuario"), req.Form.Get("Password")) {
				fmt.Println("Login ok")
				mensaje = "Usuario: " + req.Form.Get("Usuario") + ", Password: " + req.Form.Get("Password")
			} else {
				fmt.Println("Login error")
				mensaje = "Error"
			}
			response(w, true, mensaje)
		}
	default:
		response(w, false, "Comando inválido")
	}
}

/*func handler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/read" {
		fmt.Fprintf(w, r.RequestURI)
		file, err := os.Open("datos.txt")
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()
		str, err := ioutil.ReadAll(file)
		fmt.Println(string(str))
	} else if r.RequestURI == "/write" {
		fmt.Fprintf(w, r.RequestURI)
		fileW, err := os.Create("result.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer fileW.Close()
		fmt.Fprintf(fileW, "user: Pepe\nmasterKey: asdf\n")
	} else {
		fmt.Fprintf(w, "I don't know")
	}
}*/
/*
func main() {
	http.HandleFunc("/", handler)
	// Start the HTTPS server in a goroutine
	go http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
	// Start the HTTP server and redirect all incoming connections to HTTPS
	http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))
}
*/
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
