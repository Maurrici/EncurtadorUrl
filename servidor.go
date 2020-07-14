package main

import (
	"fmt"
	"github.com/Maurrici/EncurtadorUrl/url"
	"log"
	"net/http"
)

var (
	porta int
	urlBase string
)

func init(){
	porta = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", porta)
	url.ConfigurarRepositorio(url.NovoRepositorioMemoria())
}

func main(){
	http.HandleFunc("/api/encurtar", Encurtar)
	//http.HandleFunc("/r/", Redirecionar)

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
	responderCom(w,status,Headers{"Location": urlCurta})
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
