package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func main() {
	conf, err := LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	authConfig := oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://localhost:%d/callback", conf.ListenPort),
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.AuthorizeURL,
			TokenURL: conf.TokenURL,
		},
	}

	g := gin.Default()
	g.LoadHTMLGlob("templates/*")

	cookieStore := cookie.NewStore([]byte("ThisIsJustATestAppSorry"))
	g.Use(sessions.Sessions("oauthtestsession", cookieStore))

	g.GET("/", func(c *gin.Context) {
		verifier := generateToken()
		s := sessions.Default(c)
		s.Set("verifier", verifier)
		s.Save()

		verifierSum := sha256.Sum256([]byte(verifier))
		challenge := base64.URLEncoding.EncodeToString(verifierSum[:])

		url := authConfig.AuthCodeURL(
			"state",
			oauth2.AccessTypeOffline,
			oauth2.SetAuthURLParam("code_challenge", challenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		)

		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	g.GET("/callback", func(c *gin.Context) {
		ctx := context.Background()
		s := sessions.Default(c)
		verifier := s.Get("verifier").(string)

		tok, err := authConfig.Exchange(
			ctx,
			c.Query("code"),
			oauth2.SetAuthURLParam("code_verifier", verifier),
		)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		// client := authConfig.Client(ctx, tok)

		c.HTML(http.StatusOK, "token.html", tok)
	})

	g.Run(fmt.Sprintf(":%d", conf.ListenPort))
}

func generateToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
