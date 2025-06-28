package githubsync

import (
	"context"
	"fmt"
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

// func CreateIssue(ctx context.Context, client *github.Client, owner, repo, title, body string) error {
// 	issue := &github.IssueRequest{
// 		Title: &title,
// 	}
// 	_, _, err := client.Issues.Create(ctx, owner, repo, issue)
// 	return err
// }

func CreateIssueIfNotExists(ctx context.Context, client *github.Client, owner, repo, todo, body string, existingIssues map[string]struct{}) error {
	if _, found := existingIssues[todo]; found {
		return fmt.Errorf("Issue already exists: %s\n", todo)
	}

	issue := &github.IssueRequest{
		Title: &todo,
		Body:  &body,
	}

	_, _, err := client.Issues.Create(ctx, owner, repo, issue)
	return err
}

func FetchIssues(ctx context.Context, client *github.Client, owner, repo string) (map[string]struct{}, error) {
	issues := make(map[string]struct{})
	opts := &github.IssueListByRepoOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		list, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, issue := range list {
			if issue.IsPullRequest() {
				continue // Skip pull requests
			}
			issues[issue.GetTitle()] = struct{}{}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	return issues, nil
}
