package main

import (
	"fmt"
	"github.com/Maurrici/EncurtadorUrl/url"
	"log"
	"net/http"
	"strings"
	"encoding/json"
)

var (
	porta int
	urlBase string
	stats chan string
)

func init(){
	porta = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", porta)
	url.ConfigurarRepositorio(url.NovoRepositorioMemoria())
}

func main(){
	stats = make(chan string)
	defer close(stats)
	go RegistrarEstatistica(stats)

	http.HandleFunc("/api/encurtar", Encurtar)
	http.HandleFunc("/r/", Redirecionar)
	http.HandleFunc("/api/stats/", Visualizar)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",porta),nil))
}


type Headers map[string]string

func Encurtar(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		responderCom(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})

		return
	}

	u, nova, err := url.BuscarOuCriarNovaUrl(extrairUrl(r))

	if err != nil{
		responderCom(w, http.StatusBadRequest, nil)

		return
	}

	var status int
	if nova{
		status = http.StatusCreated
	}else{
		status = http.StatusOK
	}

	urlCurta := fmt.Sprintf("%s/r/%s", urlBase,u.Id)
	responderCom(w,status,
			Headers{"Location": urlCurta,
							"Link": fmt.Sprintf("%s/api/stats/%s",urlBase,u.Id)})
}

func responderCom(w http.ResponseWriter, status int, headers Headers){
	for k, v := range headers{
		w.Header().Set(k,v)
	}
	w.WriteHeader(status)
}

func extrairUrl(r *http.Request) string{
	u := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(u) //Copiando bytes da requisição
	return string(u)
}

func Redirecionar(w http.ResponseWriter, r *http.Request){
	caminho := strings.Split(r.URL.Path,"/")
	id := caminho[len(caminho)-1]

	if u := url.Buscar(id); u != nil{
		http.Redirect(w,r,u.Destino,http.StatusMovedPermanently)
		stats <- u.Id
	}else{
		http.NotFound(w,r)
	}
}

func RegistrarEstatistica(ids <-chan string){
	for id := range ids{
		url.RegistrarClick(id)
		fmt.Printf("Clique Registrado para %s",id)
	}
}

func Visualizar(w http.ResponseWriter, r *http.Request){
	caminho := strings.Split(r.URL.Path,"/")
	id := caminho[len(caminho)-1]

	if u := url.Buscar(id); u != nil {
		json, err := json.Marshal(u.Stats())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responderComJson(w, string(json))
	}else{
		http.NotFound(w,r)
	}
}

func responderComJson(w http.ResponseWriter, resposta string){
	responderCom(w, http.StatusOK, Headers{"Content-Type": "application/json"})
	fmt.Fprintf(w, resposta)
}