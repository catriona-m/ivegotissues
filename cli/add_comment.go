package cli

import (
	"fmt"

	"github.com/google/go-github/v52/github"
	c "github.com/gookit/color"
	"github.com/ivegotissues/lib/gh"
)

type AddComment struct {
	Labels     []string
	State      string
	Issues     []int
	Comment    string
	Batch      int
	OpenIssues bool
	Token      string
	Repo       string
	DryRun     bool
}

func (ac AddComment) AddComment() error {

	repo := gh.NewRepo(ac.Repo, ac.Token)

	if len(ac.Issues) > 0 {
		err := ac.addCommentToIssueList(repo)
		if err != nil {
			return err
		}

	} else if len(ac.Labels) > 0 {
		err := ac.addCommentToIssuesFilteredByLabels(repo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ac AddComment) addCommentToIssuesFilteredByLabels(repo gh.Repo) error {

	opts := github.IssueListByRepoOptions{
		State:  ac.State,
		Labels: ac.Labels,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		issues, nextPage, err := repo.ListIssuesByRepo(opts)
		if err != nil {
			return fmt.Errorf("retrieving issues from github from page %d: %v", opts.ListOptions.Page, err)
		}
		for _, issue := range issues {

			c.Printf("Adding comment to <cyan>#%d</>\t%s\t%s\n", issue.GetNumber(), issue.GetTitle(), issue.GetHTMLURL())

			if !ac.DryRun {
				err = repo.AddCommentToIssue(ac.Comment, issue.GetNumber())
				if err != nil {
					return err
				}
			}
		}

		if nextPage == 0 {
			break
		}
		opts.ListOptions.Page = nextPage

	}
	return nil
}

func (ac AddComment) addCommentToIssueList(repo gh.Repo) error {

	for _, issue := range ac.Issues {

		c.Printf("Adding comment to <cyan>#%d</>\n", issue)

		if !ac.DryRun {
			err := repo.AddCommentToIssue(ac.Comment, issue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
