package main

import (
	"crypto/aes"
	"crypto/cipher"
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
	"time"
)

type Entrada struct {
	User       string
	Password   string
	Comentario string
}

type Usuario struct {
	Sal       string
	Key       string
	MasterKey string
	Lista     map[string]Entrada
}

var userNameSession string
var sesiones = map[string]time.Time{"usuario": time.Now()}

func crearsesion(usuario string) {

	sesiones[usuario] = time.Now()

}

func comprobarsesion(usuario string) (comp bool, mensaje string) {
	tiempo_max := 60
	if val, ok := sesiones[usuario]; ok {
		if time.Now().Sub(val) < time.Duration(tiempo_max)*time.Second {
			sesiones[usuario] = time.Now()
			return true, ""
		} else {
			delete(sesiones, usuario)
			return false, "El tiempo de sesión ha expirado."
		}
	}
	fmt.Println(usuario)
	return false, "No se ha hecho login"
	/*if sesiones[usuario]!=nil && time.Now().Sub(sesiones[usuario])<time.Duration(tiempo_max)*time.Nanosecond {
		return true;
	}
	return false;*/

}

func getBaseDatos() map[string]Usuario {
	usuarios := map[string]Usuario{}
	file, err := os.Open("bd.json")
	chk(err)
	defer file.Close()
	str, err := ioutil.ReadAll(file)
	if erru := json.Unmarshal(str, &usuarios); erru != nil {
		panic(erru)
	}
	return usuarios
}

func encrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)+16)    // reservamos espacio para el IV al principio
	rand.Read(out[:16])                 // generamos el IV
	blk, err := aes.NewCipher(key)      // cifrador en bloque (AES), usa key
	chk(err)                            // comprobamos el error
	ctr := cipher.NewCTR(blk, out[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out[16:], data)    // ciframos los datos
	return
}

// función para descifrar (con AES en este caso)
func decrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)-16)     // la salida no va a tener el IV
	blk, err := aes.NewCipher(key)       // cifrador en bloque (AES), usa key
	chk(err)                             // comprobamos el error
	ctr := cipher.NewCTR(blk, data[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out, data[16:])     // desciframos (doble cifrado) los datos
	return
}

func encriptarUserPass(usuario string, usuarioSitio string, password string) (string, string) {

	keyClient := sha512.Sum512([]byte(getUserKey(usuario)))
	keyDataPass := keyClient[32:64]
	keyDataUser := keyClient[0:32]
	userEncripted := encode64(encrypt([]byte(usuarioSitio), keyDataUser))
	passEncripted := encode64(encrypt([]byte(password), keyDataPass))
	return userEncripted, passEncripted
}

func desencriptarUserPass(usuario string, usuarioSitio string, password string) (string, string) {

	keyClient := sha512.Sum512([]byte(getUserKey(usuario)))
	keyDataPass := keyClient[32:64]
	keyDataUser := keyClient[0:32]
	userDesencripted := string(decrypt(decode64(usuarioSitio), keyDataUser))
	passDesencripted := string(decrypt(decode64(password), keyDataPass))
	return userDesencripted, passDesencripted
}

func addEntry(usuario string, entrada string, usuarioSitio string, password string, comentario string) (bool, string) {
	usuarios := getBaseDatos()
	ok := false
	msg := ""

	if _, value := usuarios[usuario]; value {
		if _, entr := usuarios[usuario].Lista[entrada]; entr{
			msg = "El sitio ya existe"
		}else{
			userEncripted, passEncripted := encriptarUserPass(usuario, usuarioSitio, password)
			entrada_nueva := Entrada{User: userEncripted, Password: passEncripted, Comentario: comentario}
			usuarios[usuario].Lista[entrada] = entrada_nueva
			usuarios_json, err := json.MarshalIndent(usuarios, "", "  ")
			if err != nil {
				fmt.Println("Error marshal: ", err)
			}
			ioutil.WriteFile("bd.json", usuarios_json, 0644)
			msg = "Sitio añadido correctamente"
			ok = true
			return ok, msg
		}
	}else{
		msg = "El usuario no existe"
	}
	return ok, msg
}

func viewEntry(usuario string, entrada string) Entrada {
	usuarios := getBaseDatos()
	if value, ok := usuarios[usuario].Lista[entrada]; ok {
		
		userDesencripted, passDesencripted := desencriptarUserPass(usuario, value.User, value.Password)
		return Entrada{User:       userDesencripted,
									Password:   string(decode64(passDesencripted)),
									Comentario: value.Comentario}
	}else{
		return Entrada{User: ".", Password: ".", Comentario: "."}
	}
}

func editEntry(usuario string, entrada string, usuarioSitio string, password string, comentario string) Entrada {
	usuarios := getBaseDatos()
	userEncripted, passEncripted := encriptarUserPass(usuario, usuarioSitio, password)
	usuarios[usuario].Lista[entrada] = Entrada{
		User: userEncripted, Password: passEncripted, Comentario: comentario}

	usuarios_json, err := json.MarshalIndent(usuarios, "", "  ")
	if err != nil {
		fmt.Println("Error marshal: ", err)
	}
	ioutil.WriteFile("bd.json", usuarios_json, 0644)
	return usuarios[usuario].Lista[entrada]
}

func deleteEntry(usuario string, entrada string) bool {
	usuarios := getBaseDatos()
	if existsEntry(usuario,entrada){
			delete(usuarios[usuario].Lista, entrada)
			usuarios_json, err := json.MarshalIndent(usuarios, "", "  ")
			if err != nil {
				fmt.Println("Error marshal: ", err)
			}
			ioutil.WriteFile("bd.json", usuarios_json, 0644)
			return true
		}
	return false
}

func getUserKey(usuario string) string {
	usuarios := getBaseDatos()
	for key, value := range usuarios {
		if key == usuario {
			return value.Key
		}
	}
	return ""
}

func login(user string, password string) bool {
	usuarios := getBaseDatos()
	for key, value := range usuarios {
		if key == user {
			sha_512 := sha512.New()
			sha_512.Write([]byte(password))
			pass2 := encode64(sha_512.Sum([]byte(value.Sal)))
			if value.MasterKey == pass2 {
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

func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	chk(err)                                     // comprobamos el error
	return b                                     // devolvemos los datos originales
}

func existsEntry(usuario string, entrada string) bool {
	usuarios := getBaseDatos()
	if value, ok := usuarios[usuario].Lista[entrada]; ok {
		fmt.Println(value)
		return true
	}
	return false
}

func registro(user string, password string) bool {
	usuarios := getBaseDatos()
	entradas := map[string]Entrada{}
	if _, ok := usuarios[user]; ok {
		return false;
	}
	salG, error := GenerateRandomString(10)
	chk(error)

	sha_512 := sha512.New()
	sha_512.Write([]byte(password))
	pass2 := encode64(sha_512.Sum([]byte(salG)))
	keyG, error := GenerateRandomString(10)
	chk(error)

	//entrada_nueva := Entrada{Sitio: "", User: "", Password: "", Comentario: ""}
	//entradas[0] = entrada_nueva
	usuario_nuevo := Usuario{Sal: salG, Key: keyG, MasterKey: pass2,Lista: entradas}
	usuarios[user] = usuario_nuevo //añadimos el nuevo usuario al map
	usuarios_json, err := json.MarshalIndent(usuarios, "", "  ")
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

type respE struct {
	Ok           bool    // true -> correcto, false -> error
	ValorEntrada Entrada // Entrada
}

func response(w io.Writer, ok bool, msg string) {
	r := resp{Ok: ok, Msg: msg}    // formateamos respuesta
	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante
}

func responseEntry(w io.Writer, ok bool, entrada Entrada) {
	r := respE{Ok: ok, ValorEntrada: entrada} // formateamos respuesta
	rJSON, err := json.Marshal(&r)            // codificamos en JSON
	chk(err)                                  // comprobamos error
	w.Write(rJSON)                            // escribimos el JSON resultante
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
				crearsesion(req.Form.Get("Usuario"))
				fmt.Println("Login ok")
				mensaje = "Usuario: " + req.Form.Get("Usuario") + ", ha hecho login."
				response(w, true, mensaje)
			} else {
				fmt.Println("Login error")
				mensaje = "Usuario y/o password incorrecto"
				response(w, false, mensaje)
			}

		}
	case "Session":
		{
			comp, mens := comprobarsesion(req.Form.Get("Usuario"))
			if comp == false {
				response(w, false, mens)
			} else {
				response(w, true, "Sesión Ok")
			}
		}
	case "Add":
		{
			comp, mens := comprobarsesion(req.Form.Get("Usuario"))
			if comp == false {
				response(w, false, mens)
			} else {
				ok, msg := addEntry(req.Form.Get("Usuario"), req.Form.Get("Sitio"), req.Form.Get("Usuariositio"),
					req.Form.Get("Password"), req.Form.Get("Comentario"))
				if ok {
					response(w, ok, msg)
				} else {
					response(w, ok, msg)
				}

			}
		}
	case "View":
		{
			comp, mens := comprobarsesion(req.Form.Get("Usuario"))
			if comp == false {
				response(w, false, mens)
			} else {
				responseEntry(w, true, viewEntry(req.Form.Get("Usuario"), req.Form.Get("Sitio")))
			}
		}
	case "Delete":
		{
			comp, mens := comprobarsesion(req.Form.Get("Usuario"))
			if comp == false {
				response(w, false, mens)
			} else {
				if deleteEntry(req.Form.Get("Usuario"), req.Form.Get("Sitio")) {
					response(w, true, "Entrada eliminada")
				} else {
						response(w, false, "Entrada no eliminada")
				}

			}
		}
	case "Logout":
		{
			delete(sesiones, req.Form.Get("Usuario"))
			response(w, true, "Logout correcto")
		}
	case "Edit?":
		{
			comp, mens := comprobarsesion(req.Form.Get("Usuario"))
			if comp == false {
				response(w, false, mens)
			} else {
				response(w, existsEntry(req.Form.Get("Usuario"), req.Form.Get("Sitio")), "-")
			}
		}
	case "Edit":
		{
			responseEntry(w, true, editEntry(req.Form.Get("Usuario"), req.Form.Get("Sitio"), req.Form.Get("Usuariositio"),
				req.Form.Get("Password"), req.Form.Get("Comentario")))
		}
	default:
		response(w, false, "Comando inválido")
	}
}

func main() {

	fmt.Println("Servidor encendido en el puerto 10443...")
	crearsesion("hola")

	fmt.Println(time.Now().Sub(sesiones["hola"]))
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
	log.Println("Servidor detenido correctamente")
}
