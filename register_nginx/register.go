package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/jmcvetta/randutil"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/skratchdot/open-golang/open"

	"golang.org/x/oauth2"
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
	OIDC            ServerResponse
	Introspection   IntrospectionClient
	ResolverAddress string
}

type IntrospectionClient struct {
	ClientID           string   `json:"clientId"`
	ClientSecret       string   `json:"clientSecret"`
	Scope              []string `json:"scope"`
	AllowIntrospection bool     `json:"allowIntrospection"`
}

const registerURL = "http://localhost:8080/openid-connect-server-webapp/register"

func main() {
	addr := flag.String("resolver", "8.8.8.8", "IP of the resolver for nginx to use")
	flag.Parse()
	clientInfo := ClientInfo{RedirectURIs: []string{"http://localhost:5000/redirect_uri"},
		ClientName: "Nginx Relying Party"}

	oauthConfig := getOAuthConfig()
	state, _ := randutil.AlphaString(10)
	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	open.Run(url)
	channel := make(chan *oauth2.Token, 1)
	e := echo.New()
	e.Get("/redirect_uri", RedirectHandler(state, oauthConfig, channel))
	go e.Run(standard.New(":3333"))
	tok := <-channel
	client := oauthConfig.Client(oauth2.NoContext, tok)
	cc := &ConfigContext{OIDC: registerOAuthClient(clientInfo), ResolverAddress: *addr,
		Introspection: registerIntrospectionClient(client)}
	conf, _ := os.Create("nginx.conf")
	defer conf.Close()
	templateFile, _ := ioutil.ReadFile("nginx.conf.tmpl")
	template, _ := template.New("conf").Parse(string(templateFile))
	template.Execute(conf, cc)
}

func registerOAuthClient(info ClientInfo) ServerResponse {
	serverResponse := ServerResponse{}
	postJSON(http.DefaultClient, registerURL, &info, &serverResponse)
	return serverResponse
}

func getOAuthConfig() *oauth2.Config {
	sr := registerOAuthClient(ClientInfo{RedirectURIs: []string{"http://localhost:3333/redirect_uri"},
		ClientName: "Setup Client"})

	oauthConfig := &oauth2.Config{
		ClientID:     sr.ClientID,
		ClientSecret: sr.ClientSecret,
		Scopes:       strings.Split(sr.Scope, " "),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8080/openid-connect-server-webapp/authorize",
			TokenURL: "http://localhost:8080/openid-connect-server-webapp/token",
		},
	}

	return oauthConfig
}

func registerIntrospectionClient(client *http.Client) IntrospectionClient {
	ic := IntrospectionClient{}
	ic.AllowIntrospection = true
	ic.Scope = []string{"user/Observation.read", "user/Patient.read"}
	icResponse := IntrospectionClient{}
	postJSON(client, "http://localhost:8080/openid-connect-server-webapp/api/clients", &ic, &icResponse)
	return icResponse
}

func postJSON(client *http.Client, url string, body interface{}, response interface{}) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(body)
	resp, err := client.Post(url, "application/json", &buf)
	defer resp.Body.Close()
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Got response code: %d\n", resp.StatusCode)
		fmt.Println("Response:")
		serverResponse, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(serverResponse))
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(response)
}

func RedirectHandler(state string, config *oauth2.Config, channel chan<- *oauth2.Token) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.Query("code")
		tok, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			return err
		}
		channel <- tok
		return c.String(http.StatusOK, "Redirected properly")
	}
}
