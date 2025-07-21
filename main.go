package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alexflint/go-arg"

	"github.com/1buran/workdiary/config"
	"github.com/1buran/workdiary/internal/application/usecase"
	"github.com/1buran/workdiary/internal/infrastructure/apiclient"
)

type Date struct {
	t time.Time
}

func (d *Date) UnmarshalText(b []byte) (err error) {
	d.t, err = time.Parse(time.DateOnly, string(b))
	return
}

func (d Date) Time() time.Time { return d.t }

type TrackMode struct {
	Project  string  `arg:"-P,--project" help:"project name in config"`
	Activity string  `arg:"-A,--activity" help:"actuvity id (redmine specific)"`
	Date     Date    `arg:"-D,--date" help:"date, format: 2025-12-25"`
	Issue    string  `arg:"-I,--issue" help:"issue ID"`
	Comment  string  `arg:"-C,--comment" help:"comment"`
	Hours    float32 `arg:"-H,--hours,required" help:"hours"`
}

type CalendarMode struct {
	Month int `arg:"-m,--month" help:"choose month (default: current month)"`
	Year  int `arg:"-y,--year" help:"choose year (default: current year)"`
}

type DemoMode struct{}

type Args struct {
	Calendar *CalendarMode `arg:"subcommand:cal" help:"calendar"`
	Track    *TrackMode    `arg:"subcommand:track" help:"track time"`
	Demo     *DemoMode     `arg:"subcommand:demo" help:"demo"`

	Debug              bool   `arg:"-d,--debug" help:"enable debug"`
	Config             string `arg:"-c,--" placeholder:"CONFIG" help:"config path (default: ~/.config/workdiary/config.json)"`
	PrintConfigExample bool   `arg:"-p,--print-conf-example" help:"print example of config and exit"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	if args.PrintConfigExample {
		config.PrintDefaultConfig()
		return
	}

	if args.Demo != nil {
		usecase.Demo()
	}

	cfg, err := config.ReadConfig(args.Config)
	if err != nil {
		panic(err)
	}

	clients := make(map[string]apiclient.ApiClient)
	for _, r := range cfg.Infra.ApiClient.Redmine {
		if !r.Disabled {
			if _, ok := clients[r.Name]; ok {
				fmt.Printf("Error: project name collision, %q already parsed", r.Name)
				continue
			}
			clients[r.Name] = apiclient.NewRedmineApiClient(
				r.Name, r.Url, r.Token, r.UserId,
				r.EmployeeProfile.HourlyRate, r.LogEnabled)
		}
	}

	for _, g := range cfg.Infra.ApiClient.Gitlab {
		if !g.Disabled {
			if _, ok := clients[g.Name]; ok {
				fmt.Printf("Error: project name collision, %q already parsed", g.Name)
				continue
			}
			clients[g.Name] = apiclient.NewGitlabApiClient(
				g.Name, g.Url, g.Token, g.ProjectPath,
				g.EmployeeProfile.HourlyRate, g.LogEnabled)
		}
	}

	if args.Calendar != nil {
		today := time.Now()
		year, month := today.Year(), today.Month()

		if args.Calendar.Month > 0 {
			month = time.Month(args.Calendar.Month)
		}
		if args.Calendar.Year > 0 {
			year = args.Calendar.Year
		}

		monthbegin := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		monthend := monthbegin.AddDate(0, 1, 0)

		var clientList []apiclient.ApiClient
		for _, v := range clients {
			clientList = append(clientList, v)
		}
		usecase.Show(
			os.Stdout,
			clientList,
			monthbegin,
			monthend,
			cfg.Color("dayOff"),
			cfg.Color("workingDay"),
			cfg.Color("expectedAmount"),
			cfg.Color("infactAmount"),
			cfg.Color("summary"),
			args.Debug,
		)
	}

	if args.Track != nil {
		if client, ok := clients[args.Track.Project]; !ok {
			fmt.Printf("Error: project %q not found or disabled!\n", args.Track.Project)
			return
		} else {
			usecase.Track(
				client, args.Track.Date.Time(), args.Track.Issue, args.Track.Activity,
				args.Track.Hours, args.Track.Comment,
			)
		}
	}
}
