HSR Reminder Go
===

Go version of HSR reminder, this reminder will fetch the train information based on env variables and send the result to a webhook.

## Usage

Create `.env` and fill up the env variables.

```
cp .env.example .env
```

Execute bin file.

```
./hsr-reminder
```

## Deploy to lambda

Build binary file for lambda

```
GOARCH=amd64 GOOS=linux go build -o lambda-hsr_reminder *.go
```

Archive to zip file

```
zip -r ${PWD##*/}_$(date +%s).zip lambda-hsr_reminder
```

Upload to lambda and set `Handler` to `lambda-hsr_reminder`
