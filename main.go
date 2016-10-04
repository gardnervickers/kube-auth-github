package main

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

type TokenReviewSpec struct {
	Token string `json:"token"`
}

type TokenReviewStatus struct {
	Authenticated bool                  `json:"authenticated"`
	User          TokenReviewStatusUser `json:"user"`
}

type TokenReviewStatusUser struct {
	Username string   `json:"user"`
	Uid      string   `json:"uid"`
	Groups   []string `json:"groups"`
}

type TokenReview struct {
	Spec   TokenReviewSpec   `json:"spec"`
	Status TokenReviewStatus `json:"status"`
}

func main() {
	r := gin.Default()
	r.POST("/", Authenticate)
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func Authenticate(c *gin.Context) {
	var tokenReview TokenReview
	c.BindJSON(&tokenReview)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tokenReview.Spec.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	teams, _, err := client.Organizations.ListUserTeams(&github.ListOptions{Page: 1, PerPage: 1000})

	if err != nil {
		c.Status(500)
	}

	fmt.Println(teams)

	c.Status(400)
}
