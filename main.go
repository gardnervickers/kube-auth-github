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
	Authenticated bool                   `json:"authenticated"`
	User          *TokenReviewStatusUser `json:"user"`
}

type TokenReviewStatusUser struct {
	Username string   `json:"user"`
	Uid      string   `json:"uid"`
	Groups   []string `json:"groups"`
}

type TokenReview struct {
	ApiVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Spec       *TokenReviewSpec   `json:"spec,omitempty"`
	Status     *TokenReviewStatus `json:"status"`
}

func main() {
	r := gin.Default()
	r.POST("/", Authenticate)
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func Authenticate(c *gin.Context) {
	var tokenReview = TokenReview{}
	c.BindJSON(&tokenReview)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tokenReview.Spec.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	teams, res, err := client.Organizations.ListTeams(os.Getenv("GITHUB_ORG"), &github.ListOptions{Page: 1, PerPage: 1000})
	if err != nil {
		c.Status(res.StatusCode)
	}

	tokenReview.Status = &TokenReviewStatus{}
	tokenReview.Status.Authenticated = false

	var user *github.User
	user, res, err = client.Users.Get("")
	if err != nil {
		c.Status(res.StatusCode)
	}
	tokenReview.Status.User = &TokenReviewStatusUser{}
	tokenReview.Status.User.Username = *user.Login
	tokenReview.Status.User.Uid = fmt.Sprintf("%d", *user.ID)

	status := 403
	for _, team := range teams {
		isMember := false
		isMember, res, err = client.Organizations.IsTeamMember(*team.ID, tokenReview.Status.User.Username)
		if err != nil {
			c.Status(res.StatusCode)
		}
		if isMember {
			tokenReview.Status.User.Groups = append(tokenReview.Status.User.Groups, *team.Name)
			if *team.Name == os.Getenv("GITHUB_TEAM") {
				status = 200
				tokenReview.Spec = nil
				tokenReview.Status.Authenticated = true
			}
		}
	}

	c.JSON(status, tokenReview)
}
