# lyrical
A simple Discord bot that is capable of playing music in voice channels and running polls. This was a project I did in order to learn Golang, so if you see anything unidiomatic here it's probably because I really lacked the experience to write good Go code. Nevertheless, I did try to keep test coverage high and it seems like the features behave pretty well.

## Installation
### Dependencies
```
sudo apt-get install ffmpeg --yes
```
### Architecture specific
Build the bot for the architecture that you desire.
```
env GOOS=linux GOARCH=arm go build -o builds/discordbot .
```

### Config file
It is necessary to have a `configs/config.json` file, a sample file can be found at `configs/config.json.example`


### Run
```
./builds/discordbot
```

## Deployment
In order to facilitate rapid prototyping, i.e. painless deployment to my server, I chose to make use of a Ruby tool called Capistrano to handles deployment to a server. All it does is run SSH commands on remote servers.

In short, we can also define diffent deployment steps based on the environment. Deployment can be done as so depending on the environment.

``` bash
# deploy to the staging environment
$ bundle exec cap staging deploy

# deploy to the production environment
$ bundle exec cap production deploy
```

## How does music playback work?
Members of the server will be able to use commands to add songs 1-by-1 to the playlist of the bot. The songs will be downloaded by the bot from Youtube, and the song will added to a sqlite database. The bot will maintain a queue of songs to play based on requests from members. Once all user requested songs have been played, playback will continue by selecting random songs to play from the pool of songs previously requested by members.

## Other tools
I also wrote a CLI tool to allow adding music in bulk to the playlist of the bot, it can be found in `./cmd/bulk-music-insert`. Given an ID of a Youtube playlist, it will download a certain number of songs from it and add it to the bot's music pool, i.e. adding it to the sqlite database.