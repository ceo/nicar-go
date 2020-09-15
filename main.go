package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/janeczku/go-spinner"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Debes ingresar un dominio")
		fmt.Println("Uso: nicar <dominio>")

		os.Exit(0)
	}

	domain := os.Args[1]

	var wg sync.WaitGroup

	wg.Add(1)

	s := spinner.StartNew("Chequeando dominio: " + domain)

	go checkDomain(&wg, domain, s)

	wg.Wait()
}

func checkDomain(wg *sync.WaitGroup, domain string, s *spinner.Spinner) {
	defer wg.Done()

	domainParts := strings.Split(domain, ".")

	tld := strings.Join(domainParts[1:], ".")
	dominio := domainParts[0]

	form := url.Values{}
	form.Add("txtBuscar", dominio)
	form.Add("cmbZonas", tld)
	form.Add("btn-consultar", "Buscar")

	resp, err := http.PostForm("https://nic.ar/verificar-dominio", form)
	if err != nil {
		log.Print(err)
		return
	}

	defer resp.Body.Close()

	fmt.Println()
	if resp.StatusCode == http.StatusOK {

		s.Stop()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		if strings.Contains(bodyString, "El dominio está disponible") {
			fmt.Println("Dominio disponible")
		} else if strings.Contains(bodyString, "El dominio no está disponible") {
			fmt.Println("El dominio no está disponible")
		} else if strings.Contains(bodyString, "El nombre de dominio que ingresaste no es válido.") {
			fmt.Println("Dominio invalido/reservado/no disponible.")
		} else {
			fmt.Println("No se pudo leer la respuesta correctamente.")

			// Saves error html into a file for debug purposes
			err = ioutil.WriteFile("last_nic_error.html", []byte(bodyString), 0644)
			if err != nil {
				panic(err)
			}
		}

	}

	return
}
