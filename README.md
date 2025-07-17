# workdiary
Simple self-motivation cli app for tracking work time in Redmine, GitLab,
show calendar with month earnings.
![Main demo](https://i.imgur.com/M1wBOlm.gif)

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
          "Disabled": false,
          "LogEnabled": true,
          "employee": {
            "DailyHoursLimit": 8,
            "HourlyRate": 10.5
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

## Tasks

These are tasks of [xc](https://github.com/joerdav/xc) runner.

### vhs

Run VHS fo update gifs.

```
vhs demo/main.tape
```

### imgur

Upload to Imgur and update readme.

```
declare -A demo=()
demo["main"]="Main demo"

for i in ${!demo[@]}; do
    . .env && url=`curl --location https://api.imgur.com/3/image \
        --header "Authorization: Client-ID ${clientId}" \
        --form image=@demo/$i.gif \
        --form type=image \
        --form title=workdiary \
        --form description=Demo | jq -r '.data.link'`
    sed -i "s#^\!\[${demo[$i]}\].*#![${demo[$i]}]($url)#" README.md
done
```
