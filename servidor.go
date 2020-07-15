package main

import (
	"fmt"
	"github.com/Maurrici/EncurtadorUrl/url"
	"log"
	"net/http"
	"strings"
	"encoding/json"
	"flag"
)

var (
	porta *int
	logLigado *bool
	urlBase string
)

func init(){
	porta = flag.Int("p", 8888, "porta")
	logLigado = flag.Bool("l", true, "log ligado/desligado")

	flag.Parse()

	urlBase = fmt.Sprintf("http://localhost:%d", *porta)
	url.ConfigurarRepositorio(url.NovoRepositorioMemoria())
}

func main(){
	stats := make(chan string)
	defer close(stats)
	go RegistrarEstatistica(stats)

	http.HandleFunc("/api/encurtar", Encurtar)
	http.Handle("/r/", &Redirecionar{stats})
	http.HandleFunc("/api/stats/", Visualizar)

	logar("Iniciando Servidor na porta %d...",*porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",*porta),nil))
}

type Redirecionar struct {
	stats chan string
}

func (red *Redirecionar) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	buscarUrlEExecutar(w, r, func(u *url.Url) {
		http.Redirect(w, r, u.Destino, http.StatusMovedPermanently)
		red.stats <- u.Id
	})
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

	logar("URL %s encurtada com sucesso para %s",u.Destino,urlCurta)

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

func RegistrarEstatistica(ids <-chan string){
	for id := range ids{
		url.RegistrarClick(id)
		logar("Clique Registrado com Sucesso para %s",id)
	}
}

func Visualizar(w http.ResponseWriter, r *http.Request){
	buscarUrlEExecutar(w, r, func(u *url.Url) {
		json, err := json.Marshal(u.Stats())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		responderComJson(w,string(json))
	})
}

func responderComJson(w http.ResponseWriter, resposta string){
	responderCom(w, http.StatusOK, Headers{"Content-Type": "application/json"})
	fmt.Fprintf(w, resposta)
}

func buscarUrlEExecutar(w http.ResponseWriter, r *http.Request, executor func(u *url.Url)){
	caminho := strings.Split(r.URL.Path,"/")
	id := caminho[len(caminho)-1]

	if u := url.Buscar(id); u != nil{
		executor(u)
	}else{
		http.NotFound(w, r)
	}
}

func logar(formato string, valores ...interface{}){
	if *logLigado{
		log.Printf(fmt.Sprintf("%s\n",formato),valores...)
	}
}