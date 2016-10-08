package main

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

func main() {
	validateEnvironment("PORT", "GITHUB_ORG", "GITHUB_TEAM")
	r := gin.Default()
	r.POST("/", authenticate)
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func validateEnvironment(varNames ...string) {
	valid := true
	for _, varName := range varNames {
		if os.Getenv(varName) == "" {
			fmt.Fprintf(os.Stdout, "invalid environment: missing $%v\n", varName)
			valid = false
		}
	}
	if !valid {
		os.Exit(1)
	}
}

func authenticate(c *gin.Context) {
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
