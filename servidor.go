package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"time"

	"github.com/Maurrici/EncurtadorUrl/url"
)

var (
	porta int
	urlBase string
)

func init(){
	porta = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", porta)
}

func main(){
	http.HandleFunc("/api/encurtar", Encurtar)
	http.HandleFunc("/r/", Redirecionar)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",porta),nil))
}


type Headers map[string]string

func Encurtar(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		responderCom(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})

		return
	}

	url, nova, err := url.BuscarOuCriarNovaUrl(extrairUrl(r))

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

	urlCurta := fmt.Sprintf("%s/r/%s", urlBase,url.Id)
	responderCom(w,status,Headers{"Location": urlCurta})
}

func responderCom(w http.ResponseWriter, status int, headers Headers){
	for k, v := range headers{
		w.Header().Set(k,v)
	}
	w.WriteHeader(status)
}

func extrairUrl(r *http.Request) string{
	url := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(url)//Copiando bytes da requisição
	return string(url)
}
