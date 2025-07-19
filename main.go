package main

import (
	"os"
	"time"

	"github.com/alexflint/go-arg"

	"github.com/1buran/workdiary/config"
	"github.com/1buran/workdiary/internal/application/usecase"
	"github.com/1buran/workdiary/internal/infrastructure/apiclient"
)

type CalendarMode struct {
	Month int `arg:"-m,--month" help:"choose month (default: current month)"`
	Year  int `arg:"-y,--year" help:"choose year (default: current year)"`
}

type Args struct {
	Calendar *CalendarMode `arg:"subcommand:cal" help:"calendar"`

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

	cfg, err := config.ReadConfig(args.Config)
	if err != nil {
		panic(err)
	}

	if args.Calendar != nil {
		var clients []apiclient.ApiClient
		for _, r := range cfg.Infra.ApiClient.Redmine {
			if !r.Disabled {
				clients = append(
					clients, apiclient.NewRedmineApiClient(
						r.Name, r.Url, r.Token, r.UserId,
						r.EmployeeProfile.HourlyRate, r.LogEnabled))
			}
		}

		for _, g := range cfg.Infra.ApiClient.Gitlab {
			if !g.Disabled {
				clients = append(
					clients, apiclient.NewGitlabApiClient(
						g.Url, g.Token, g.ProjectPath,
						g.EmployeeProfile.HourlyRate, g.LogEnabled))
			}
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
		monthend := monthbegin.AddDate(0, 1, 0)

		usecase.Show(
			os.Stdout,
			clients,
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
}
