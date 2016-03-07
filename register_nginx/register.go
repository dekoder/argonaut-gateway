package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

type ClientInfo struct {
	RedirectURIs []string `json:"redirect_uris"`
	ClientName   string   `json:"client_name"`
}

type ServerResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
}

type ConfigContext struct {
	ServerResponse
	ResolverAddress string
}

const registerURL = "http://localhost:8080/openid-connect-server-webapp/register"

func main() {
	addr := flag.String("resolver", "8.8.8.8", "IP of the resolver for nginx to use")
	flag.Parse()
	clientInfo := &ClientInfo{RedirectURIs: []string{"http://localhost:5000/redirect_uri"},
		ClientName: "Nginx Relying Party"}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(clientInfo)
	resp, err := http.Post(registerURL, "application/json", &buf)
	defer resp.Body.Close()
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode == http.StatusCreated {
		templateFile, _ := ioutil.ReadFile("nginx.conf.tmpl")
		template, _ := template.New("conf").Parse(string(templateFile))
		decoder := json.NewDecoder(resp.Body)
		serverResponse := ServerResponse{}
		decoder.Decode(&serverResponse)
		cc := &ConfigContext{ServerResponse: serverResponse, ResolverAddress: *addr}
		conf, _ := os.Create("nginx.conf")
		defer conf.Close()
		template.Execute(conf, cc)
	} else {
		fmt.Printf("Got response code: %d\n", resp.StatusCode)
		fmt.Println("Response:")
		serverResponse, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(serverResponse)
	}
}
