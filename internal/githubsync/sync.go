package githubsync

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v72/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

// TODO Add label
func CreateIssueIfNotExists(ctx context.Context, client *github.Client, owner, repo, todoTitle, body string, existingIssues map[string]*github.Issue) error {
	if _, found := existingIssues[todoTitle]; found {
		return fmt.Errorf("Issue already exists: %s\n", todoTitle)
	}

	issue := &github.IssueRequest{
		Title: &todoTitle,
		Body:  &body,
	}

	_, _, err := client.Issues.Create(ctx, owner, repo, issue)
	return err
}

func FetchAllOpenIssues(ctx context.Context, client *github.Client, owner, repo string) (map[string]*github.Issue, error) {
	issues := make(map[string]*github.Issue)
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
			issues[issue.GetTitle()] = issue
		}

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	return issues, nil
}

func UpdateIssueIfNeeded(ctx context.Context, client *github.Client, owner, repo string, issue *github.Issue, updateField string, update []string) error {
	if issue == nil {
		return fmt.Errorf("issue is nil")
	}

	switch updateField {
	case "labels":
		currentLabels := make(map[string]bool)
		for _, lbl := range issue.Labels {
			currentLabels[lbl.GetName()] = true
		}

		var toAdd []string
		for _, lbl := range update {
			if !currentLabels[lbl] {
				toAdd = append(toAdd, lbl)
			}
		}

		if len(toAdd) == 0 {
			return nil
		}
		_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, issue.GetNumber(), toAdd)
		return err
	case "body":
		body := strings.Join(update, "\n")
		issueRequest := &github.IssueRequest{
			Body: &body,
		}
		_, _, err := client.Issues.Edit(ctx, owner, repo, issue.GetNumber(), issueRequest)
		return err
	default:
		return fmt.Errorf("unsupported update field: %s", updateField)
	}
}
