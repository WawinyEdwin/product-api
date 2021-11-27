package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/WawinyEdwin/product-api/working/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

//ServeHTTP is the maiin entry point for the handler
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//create products operation
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
	}

	//create products operation
	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)

		//expect the id in the url
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("invalid uri more than 1 id")
			http.Error(rw, "invalid uri", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("invalid uri more than one category group")
			http.Error(rw, "invalid uri", http.StatusBadRequest)
			return
		}

		idString := g[0][1]

		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("unable to convert")
			http.Error(rw, "invalid uri", http.StatusBadRequest)
			return
		}

		// validate url
		// p.l.Println("got id", id)

		p.updateProducts(id, rw, r)
		return

	}

	//catch nil
	//if no method is returned
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

//get products operation
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

//add products operation
func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to umarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

//create products operation
func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to umarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
