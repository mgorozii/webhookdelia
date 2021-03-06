# Webhookdelia

![](https://github.com/mgorozii/webhookdelia/workflows/common/badge.svg)

Webhookdelia is a bot, implementing slack-like webhooks for telegram chats. Just start conversation with it or add it to the group chat and it'll send the example URL for sending messages. You could use [@webhookdelia](http://t.me/webhookdelia_bot) or setup it on your server. 
The webhook URL is looks like `https://webhookdelia.mgorozii.com/send/<some secret uuid>`. Text may be either raw text or [html](https://core.telegram.org/bots/api#html-style). 

## Usage

In short, the endpoint supports as GET as POST requests. There are two expected params (both of them could be as query-param, as keys in POST body):
- `text` - any string
- `parse_mode` - optional, [telegram parse mode](https://core.telegram.org/bots/api#formatting-options), "Markdown" by default, must be one of "Markdown", "MarkdownV2", "HTML"

If it's POST request, `text` will be parsed as [golang html template](https://golang.org/pkg/html/template/) and executed with data from POST body.

Sending some `text` with markdown formatting:
`http https://webhookdelia.mgorozii.com/send/<replace with secret key>?text="[link to repo](https://github.com/mgorozii/webhookdelia)"`

Specify other formatting:
`http https://webhookdelia.mgorozii.com/send/<replace with secret key>?text=<b>hi!</b>&parse_mode=HTML`

Using golang templates and POST:
`http  https://webhookdelia.mgorozii.com/send/<replace with secret key> text="Hi, {{.name}}" name="Webhookdelia"`

## Building

Webhookdelia supports multiple storages (thanks to [philippgille/gokv](https://github.com/philippgille/gokv)). By default, it's file system storage, also it's possible to store the data in the PostgreSQL and Redis.
To build it with support of the required storage you need to use proper tag:
```
go build cmd/webhookdelia/main.go  # for files
go build -tags postgres cmd/webhookdelia/main.go
go build -tags redis cmd/webhookdelia/main.go
```
The similar approach for building Docker-image:
```
docker build --build-arg store=redis -t webhookdelia . # for redis
```

## Configuration

Webhookdelia uses [spf13/viper](https://github.com/spf13/viper) for configuration, so you can provide it as in the file conf.yaml (in the working dir) or as environment variables. 

### Common


| Name           | Type   | Default                  | Comment                                 |
| -------------- | ------ | ------------------------ | --------------------------------------- |
| telegram_token | string |                          | Required.                               |
| public_url     | string |                          | Required. Format: http://localhost:8080 |
| port           | string | 8080                     |                                         |
| telegram_url   | string | https://api.telegram.org |                                         |
| release        | bool   | false                    |                                         |


### Store configuration

#### File

| Name                 | Type   | Default           | Comment |
| -------------------- | ------ | ----------------- | ------- |
| store_file_directory | string | /tmp/webhookdelia |         |


#### PostgreSQL

| Name                                | Type   | Default                                           | Comment                                                                                    |
| ----------------------------------- | ------ | ------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| store_postgres_connection_url       | string | postgres://postgres@/webhookdelia?sslmode=disable | Format: postgres://username[:password]@[address]/dbname[?param1=value1&...&paramN=valueN]. |
| store_postgres_table_name           | string | webhookdelia                                      |                                                                                            |
| store_postgres_max_open_connections | int    | 100                                               |                                                                                            |

#### Redis

| Name                 | Type   | Default        | Comment |
| -------------------- | ------ | -------------- | ------- |
| store_redis_address  | string | localhost:6379 |         |
| store_redis_password | string |                |         |
| store_redis_db       | int    | 0              |         |
