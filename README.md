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

## Deploy

Install aws-cdk CLI

```
yarn global add aws-cdk
```

Initialize stack

```
yarn cdk bootstrap
```

Run deploy script

```
yarn deploy
```
