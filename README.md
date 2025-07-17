# workdiary
Simple self-motivation cli app for tracking work time in Redmine, GitLab,
show calendar with month earnings.
![Main demo](https://i.imgur.com/00ZItje.mp4)

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
        --form image=@demo/$i.webm \
        --form type=image \
        --form title=workdiary \
        --form description=Demo | jq -r '.data.link'`
    sed -i "s#^\!\[${demo[$i]}\].*#![${demo[$i]}]($url)#" README.md
done
```
