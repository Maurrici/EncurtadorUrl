package url

import "time"

type Url struct {
	Id string
	Criacao time.Time
	Destino string
}

func (url Url) BuscarOuCriarNovaUrl(urlRequest string) (Url, bool, error){

}