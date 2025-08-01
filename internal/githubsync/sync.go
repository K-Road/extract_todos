package githubsync

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/K-Road/extract_todos/internal/helper"
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
	labels := []string{"todo", "sync-generated"}
	issue := &github.IssueRequest{
		Title:  &todoTitle,
		Body:   &body,
		Labels: &labels,
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

// List github issues that aren't in the list of todos
// Interate over the issues to check if they have a sync-generated label
// close issue if it exists and it is not in the list of todos
func CloseDeletedTodos(ctx context.Context, client *github.Client, owner, repo, project string, todos map[string]struct{}) error {
	existingIssues, err := FetchAllOpenIssues(ctx, client, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to fetch existing issues: %w", err)
	}

	for title, issue := range existingIssues {
		if _, exists := todos[title]; !exists &&
			issue.GetState() == "open" &&
			hasLabel(issue.Labels, "sync-generated") {
			err := CloseIssueIfExists(ctx, client, owner, repo, issue.GetNumber())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CloseIssueIfExists(ctx context.Context, client *github.Client, owner, repo string, issueNumber int) error {
	issueRequest := &github.IssueRequest{
		State: helper.String("closed"),
	}

	_, _, err := client.Issues.Edit(ctx, owner, repo, issueNumber, issueRequest)
	if err != nil {
		return fmt.Errorf("failed to close issue #%d: %w", issueNumber, err)
	}
	log.Printf("Closed issue #%d\n", issueNumber)
	return nil
}

func hasLabel(labels []*github.Label, target string) bool {
	for _, label := range labels {
		if label.GetName() == target {
			return true
		}
	}
	return false
}
