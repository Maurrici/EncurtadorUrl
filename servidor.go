package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"

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
	http.HandleFunc("/api/encurtar", Encurtador)
	http.HandleFunc("/r/", Redirecionar)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",porta),nil))
}

