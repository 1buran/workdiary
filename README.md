# workdiary
Simple self-motivation cli app for tracking work time in Redmine, GitLab,
show calendar with month earnings.
![Main](https://i.imgur.com/iVcAzGs.mp4)
![Themes](https://i.imgur.com/oLoncRQ.png)
![Debug](https://i.imgur.com/twsBK4k.png)

## Installation

```sh
go install github.com/1buran/workdiary@latest
```

## Configuration

You can use your own color theme and change other default settings via config.

To print example of config use `-p` option:

```json
{
  "Infra": {
    "ApiClient": {
      "Redmine": [
        {
          "Name": "example",
          "Url": "http://example.com",
          "Token": "xxxxxxxxxxxxxx",
          "UserId": "100",
          "Project": "redmine-project",
          "Disabled": false,
          "LogEnabled": true,
          "employee": {
            "DailyHoursLimit": 8,
            "HourlyRate": 10.5
          }
        }
      ],
      "Gitlab": [
        {
          "Name": "example2",
          "Url": "https://domain.com/api/graphql",
          "Token": "xxxxxxxxxxxxxx",
          "ProjectPath": "group/project",
          "Disabled": false,
          "LogEnabled": true,
          "employee": {
            "DailyHoursLimit": 4,
            "HourlyRate": 15
          }
        }
      ]
    }
  },
  "App": {
    "Theme": {
      "Colors": {
        "dayOff": "#a958ad",
        "expectedAmount": "#ff9ff3",
        "infactAmount": "#4cd137",
        "summary": "#fd79a8",
        "workingDay": "#0d420d"
      }
    }
  }
}
```

Save default config somewhere and redact:
- DO NOT FORGET: change redmine/gitlab url, token and other settings to correct.
- DO NOT FORGET: correct hourly rate value to your own, look up `HourlyRate` option
- you may want to change colors: `App.Theme.Colors` section
- you may want to disable log messages: `Infra.ApiClient.Redmine.LogEnabled` option
- you may want to add more than one tracking system, just add them to appropriate
  sections: `Infra.ApiClient.Redmine` or `Infra.ApiClient.Gitlab`

## Usage

Use `-h` or `--help` to print help:

```
Usage: workdiary [-c CONFIG] [--debug] [--print-conf-example] <command> [<args>]

Options:
  -c CONFIG              config path (default: ~/.config/workdiary/config.json)
  --debug, -d            enable debug
  --print-conf-example, -p
                         print example of config and exit
  --help, -h             display this help and exit

Commands:
  cal                    calendar
```

Use `-c` to specify custom config, the workdiary uses `~/.config/workdiary/config.json`
by default.

Use `-d` or `--debug` to enable debug, in this mode workdiary print extra info about
expected working hours in month, hourly rate etc.

Use `-p` or `--print-conf-example` to print example of config.

Subcommands:
- `cal`- show calendar with stats of working days.
- `track`- register spent time in tracking system (Redmine or GitLab)

### Calendar

Use `-h` or `--help` to print help of calendar subcommand:

```
Usage: workdiary cal [--month MONTH] [--year YEAR]

Options:
  --month MONTH, -m MONTH
                         choose month (default: current month)
  --year YEAR, -y YEAR   choose year (default: current year)

Global options:
  --debug, -d            enable debug
  -c CONFIG              config path (default: ~/.config/workdiary/config.json)
  --print-conf-example, -p
                         print example of config and exit
  --help, -h             display this help and exit
```

Use `-y` or `--year` to set custom year, default: current year.

Use `-m` or `--month` to set custom month, default: current month.

### Track

Use `-h` or `--help` to print help of tracking time subcommand:

```
Usage: workdiary track [--project PROJECT] [--activity ACTIVITY] [--date DATE] [--issue ISSUE] [--comment COMMENT] --hours HOURS

Options:
  --project PROJECT, -P PROJECT
                         project name in config
  --activity ACTIVITY, -A ACTIVITY
                         actuvity id (redmine specific)
  --date DATE, -D DATE   date
  --issue ISSUE, -I ISSUE
                         issue ID
  --comment COMMENT, -C COMMENT
                         comment
  --hours HOURS, -H HOURS
                         hours

Global options:
  --debug, -d            enable debug
  -c CONFIG              config path (default: ~/.config/workdiary/config.json)
  --print-conf-example, -p
                         print example of config and exit
  --help, -h             display this help and exit
```

Use `-P` (capital letter) or `--project` to specify the project,
this is the value must be matched to a value from confi  sections `Redmine.Name` or
`Gitlab.Name`, if you need register time to more than one project,
just add another one seciton with that project to config.

Use `-A` (capital letter) or `--activity` (Redmine specific) to specify activity type.

Use `-D` (capital letter) or `--date` to specify the date of spent time,
    format: 2025-12-27, default: today.

Use `-I` (capital letter) or `--issue` to specify the issue ID inside tracking system
(required for Gitlab but optional for Redmine).

Use `-H` (capital letter) or `--hours` to specify spent time in hours (required).

Use `-C` (capital letter) or `--comment` to specify a comment (optional).

## Tasks

These are tasks of [xc](https://github.com/joerdav/xc) runner.

### vhs

Run VHS for update gifs.

```
vhs demo/main.tape
```

### themes

Show available color themes.

```
vhs demo/themes.tape
```

### debug

Debug calendar data.

```
vhs demo/debug.tape
```

### imgur

Upload to Imgur and update readme.

```
declare -A demo=()
demo["main.webm"]="Main"
demo["themes.png"]="Themes"
demo["debug.png"]="Debug"

for i in ${!demo[@]}; do
    . .env && url=`curl --location https://api.imgur.com/3/image \
        --header "Authorization: Client-ID ${clientId}" \
        --form image=@demo/$i \
        --form type=file \
        --form title=workdiary \
        --form description=Demo | jq -r '.data.link'`
    sed -i "s#^\!\[${demo[$i]}\].*#![${demo[$i]}]($url)#" README.md
done
```
