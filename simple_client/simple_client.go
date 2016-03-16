package main

import (
	"flag"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"golang.org/x/oauth2"
)

type Session struct {
	Token  *oauth2.Token
	Config *oauth2.Config
}

func main() {
	clientId := flag.String("id", "", "the Client ID for this OAuth2 Client")
	clientSecret := flag.String("secret", "", "the Client secret for this OAuth2 Client")
	flag.Parse()

	e := echo.New()

	session := &Session{}
	session.Config = &oauth2.Config{
		ClientID:     *clientId,
		ClientSecret: *clientSecret,
		Scopes:       []string{"user/Observation.read", "user/Patient.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8080/openid-connect-server-webapp/authorize",
			TokenURL: "http://localhost:8080/openid-connect-server-webapp/token",
		},
	}

	e.Use(middleware.Logger())
	e.Get("/", IndexHandler(session))
	e.Get("/redirect", RedirectHandler(session))
	e.Run(standard.New(":3333"))
}

func IndexHandler(session *Session) echo.HandlerFunc {
	return func(c echo.Context) error {
		if session.Token == nil {
			url := session.Config.AuthCodeURL("state", oauth2.AccessTypeOnline)
			return c.Redirect(http.StatusTemporaryRedirect, url)
		}
		client := session.Config.Client(oauth2.NoContext, session.Token)
		resp, err := client.Get("http://localhost:5000/api/Observation")
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		_, err = io.Copy(c.Response(), resp.Body)
		return err
	}
}

func RedirectHandler(session *Session) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.Query("code")
		tok, err := session.Config.Exchange(oauth2.NoContext, code)
		if err != nil {
			return err
		}
		session.Token = tok
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
