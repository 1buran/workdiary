package apiclient

import (
	"context"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"

	"github.com/1buran/workdiary/internal/domain/valueobject"
)

func NewGitlabApiClient(
	projectName, url, token, projectPath string,
	hourlyRate float32,
	logEnabled bool,
) ApiClient {
	grc := graphql.NewClient(url, http.DefaultClient).
		WithRequestModifier(func(r *http.Request) {
			r.Header.Set("Authorization", "Bearer "+token)
		}).
		WithDebug(logEnabled)
	return gitlabApiClient{
		client: grc, projectPath: projectPath, projectName: projectName, hourlyRate: hourlyRate,
	}
}

type gitlabApiClient struct {
	client     *graphql.Client
	hourlyRate float32 // employee rate

	projectName, projectPath string
}

func (g gitlabApiClient) List(d1, d2 time.Time) (<-chan valueobject.Day, <-chan error) {
	ch := make(chan valueobject.Day)
	er := make(chan error)

	go func() {
		defer close(ch)
		defer close(er)

		var q struct {
			Project struct {
				Issues struct {
					Nodes []struct {
						IssueID  string `graphql:"iid"`
						Title    string `graphql:"title"`
						WebUrl   string `graphql:"webUrl"`
						TimeLogs struct {
							Nodes []struct {
								TimeSpentSeconds int       `graphql:"timeSpent"`
								SpentAt          time.Time `graphql:"spentAt"`
								Summary          string    `graphql:"summary"`
								User             struct {
									Username string `graphql:"username"`
									Name     string `graphql:"name"`
								} `graphql:"user"`
							} `graphql:"nodes"`
							PageInfo struct {
								EndCursor   string `graphql:"endCursor"`
								HasNextPage bool   `graphql:"hasNextPage"`
							} `graphql:"pageInfo"`
						} `graphql:"timelogs(first: $first, after: $endCursor)"` // todo: pagination
					}
				} `graphql:"issues(first: $first)"`
			} `graphql:"project(fullPath: $projectPath)"`
		}

		variables := map[string]any{
			"first":       100,
			"projectPath": graphql.ID(g.projectPath),
			"endCursor":   "",
		}

		err := g.client.Query(context.Background(), &q, variables)
		if err != nil {
			er <- err
		}

		for _, issue := range q.Project.Issues.Nodes {
			for _, timelog := range issue.TimeLogs.Nodes {
				if timelog.SpentAt.After(d1) && timelog.SpentAt.Before(d2) {
					d := valueobject.NewDay(timelog.SpentAt)
					d.Track(g.hourlyRate, float32(timelog.TimeSpentSeconds/3600))
					ch <- d
				}
			}
		}
	}()
	return ch, er
}

func (g gitlabApiClient) Project() (p string) { return g.projectName }

// Custom Gtilab graphql type IssuableID is a global ID. It is encoded as a string.
// An example IssuableID is: "gid://gitlab/Issuable/1".
// https://docs.gitlab.com/api/graphql/reference/#issuableid
type IssuableID string

// Custom Gitlab graphQL type Time represented in ISO 8601.
// For example: “2021-03-09T14:58:50+00:00”.
// https://docs.gitlab.com/api/graphql/reference/#time
type Time string

func (g gitlabApiClient) Track(
	date time.Time,
	issue, activity string,
	hours float32,
	comment string,
) error {
	if date.IsZero() {
		date = time.Now()
	}

	var q struct {
		Project struct {
			Issue struct {
				ID    string `graphql:"iid"`
				GID   string `graphql:"id"` // Global Gitlab ID, used later as IssueableID
				Title string `graphql:"title"`
			} `graphql:"issue(iid: $issue)"`
		} `graphql:"project(fullPath: $projectPath)"`
	}

	variables := map[string]any{
		"projectPath": graphql.ID(g.projectPath),
		"issue":       issue,
	}

	if err := g.client.Query(context.Background(), &q, variables); err != nil {
		return err
	}

	var m struct {
		TimelogCreate struct {
			Timelog struct {
				ID               string    `graphql:"id"`
				TimeSpentSeconds int       `graphql:"timeSpent"`
				Comment          string    `graphql:"summary"`
				SpentAt          time.Time `graphql:"spentAt"`
				User             struct {
					Username string `graphql:"username"`
				} `graphql:"user"`
			} `graphql:"timelog"`
		} `graphql:"timelogCreate(input: {issuableId: $issue, timeSpent: $timeSpent, spentAt: $ts, summary: $comment})"`
	}

	variables = map[string]any{
		"issue":     IssuableID(q.Project.Issue.GID),
		"ts":        Time(date.Format(time.RFC3339)),
		"timeSpent": (time.Minute * time.Duration(hours*60)).String(),
		"comment":   comment,
	}

	if err := g.client.Mutate(context.Background(), &m, variables); err != nil {
		return err
	}

	return nil
}
