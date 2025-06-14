package githubsync

import (
	"context"
	"os"

	"github.com/google/go-github/v72/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func CreateIssue(client *github.Client, owner, repo, title, body string) error {
	issue := &github.IssueRequest{
		Title: &title,
	}
	_, _, err := client.Issues.Create(context.Background(), owner, repo, issue)
	return err
}
