package url

type repositorioMemoria struct {
	urls map[string]*Url
	clicks map[string]int
}

func NovoRepositorioMemoria() *repositorioMemoria{
	return &repositorioMemoria{make(map[string]*Url),make(map[string]int)}
}

func (repositorio *repositorioMemoria) IdExiste(id string) bool{
	_, existe := repositorio.urls[id]

	return existe
}

func (repositorio *repositorioMemoria) BuscaPorId(id string) *Url {
	return repositorio.urls[id]
}

func (repositorio *repositorioMemoria) BuscaPorUrl(urlRequest string) *Url{
	for _, u := range repositorio.urls{
		if u.Destino == urlRequest {
			return u
		}
	}

	return nil
}

func (repositorio *repositorioMemoria) Salvar(url Url) error{
	repositorio.urls[url.Id] = &url

	return nil
}

func (repositorio *repositorioMemoria) RegistrarClick(id string){
	repositorio.clicks[id] += 1
}

func (repositorio *repositorioMemoria) BuscarClick(id string) int{
	return repositorio.clicks[id]
}