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

// todo
func (g gitlabApiClient) Track(date time.Time, hours float32) error { return nil }
