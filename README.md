Install dependencies
```
sudo apt-get install ffmpeg --yes
go get github.com/jonas747/dca/cmd/dca
```

Build the bot for the architecture that you desire.
```
env GOOS=linux GOARCH=arm go build -o builds/discordbot .
```

Run it
```
./builds/discordbot
```