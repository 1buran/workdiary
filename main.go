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
	Month      int    `arg:"-m,--month" help:"choose month (default: current month)"`
	Year       int    `arg:"-y,--year" help:"choose year (default: current year)"`
	Theme      string `arg:"-t,--theme" help:"color theme name"`
	ListThemes bool   `arg:"-l,--list-themes" help:"lsit available color themes"`
}

type DemoMode struct {
	Month int    `arg:"-m,--month" help:"choose month (default: current month)"`
	Theme string `arg:"-t,--theme" help:"color theme name"`
}

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
		var themes [][3]string
		if len(args.Demo.Theme) > 0 {
			th := config.Themes.Get(args.Demo.Theme)
			if th == nil {
				fmt.Printf("Error: theme %q not found", args.Demo.Theme)
			}
			themes = append(
				themes, [3]string{args.Demo.Theme, th.Color("workingDay"), th.Color("dayOff")})
		} else {
			for _, name := range config.Themes.List() {
				themes = append(
					themes, [3]string{name, config.Themes.Get(name).Color("workingDay"),
						config.Themes.Get(name).Color("dayOff")})
			}
		}
		if args.Demo.Month > 0 {
			usecase.Demo2(args.Demo.Month, themes[0][1], themes[0][2])
			return
		}
		usecase.Demo(themes)
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
		if args.Calendar.ListThemes {
			config.Themes.PrintList()
			return
		}

		today := time.Now()
		year, month := today.Year(), today.Month()

		if args.Calendar.Month > 0 {
			month = time.Month(args.Calendar.Month)
		}
		if args.Calendar.Year > 0 {
			year = args.Calendar.Year
		}

		monthbegin := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		monthend := monthbegin.AddDate(0, 1, 0).Add(-time.Nanosecond)

		var clientList []apiclient.ApiClient
		for _, v := range clients {
			clientList = append(clientList, v)
		}

		doff, wday, exp, infact, sum := cfg.Color("dayOff"),
			cfg.Color("workingDay"),
			cfg.Color("expectedAmount"),
			cfg.Color("infactAmount"),
			cfg.Color("summary")

		if len(args.Calendar.Theme) > 0 {
			if t := config.Themes.Get(args.Calendar.Theme); t != nil {
				doff, wday, exp, infact, sum = t.Color("dayOff"),
					t.Color("workingDay"),
					t.Color("expectedAmount"),
					t.Color("infactAmount"),
					t.Color("summary")
			}
		}
		usecase.Show(
			os.Stdout,
			clientList,
			monthbegin,
			monthend,
			doff, wday, exp, infact, sum,
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
