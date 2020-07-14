package url

import (
	"time"
	"net/url"
	"math/rand"
)

const (
	tamanho = 5
	simbolos = "abcdefghijklmnopqr...STUVWXYZ123456789_-+"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

type Url struct {
	Id string `json:"id"`
	Criacao time.Time `json:"criacao"`
	Destino string `json:"destino"`
}

func (Url *Url) Stats() *Stats{
	clicks := repo.BuscarClick(Url.Id)
	return &Stats{Url, clicks}
}

type Stats struct {
	Url *Url `json:"url"`
	Clicks int `json:"clicks"`
}

type Repositorio interface {
	IdExiste(id string) bool
	BuscaPorId(id string) *Url
	BuscaPorUrl(url string) *Url
	Salvar(url Url) error
	RegistrarClick(id string)
	BuscarClick(id string) int
}

var repo Repositorio

func ConfigurarRepositorio(r Repositorio){
	repo = r
}

func BuscarOuCriarNovaUrl(urlRequest string) (*Url,bool,error){
	if u := repo.BuscaPorUrl(urlRequest); u != nil {
		return u, false, nil
	}

	if _, err := url.ParseRequestURI(urlRequest); err != nil {
		return nil, false, err
	}

	u := Url{gerarId(),time.Now(), urlRequest}
	repo.Salvar(u)

	return &u, true, nil
}

func gerarId() string{
	novoId := func() string{
		id := make([]byte, tamanho,tamanho)
		for i := range id{
			id[i] = simbolos[rand.Intn(len(simbolos))]
		}

		return string(id)
	}

	for{
		if id := novoId(); !repo.IdExiste(id){
			return id
		}
	}
}

func Buscar(id string) *Url{
	return repo.BuscaPorId(id)
}

func RegistrarClick(id string){
	repo.RegistrarClick(id)
}
