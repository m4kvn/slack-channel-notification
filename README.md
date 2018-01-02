# slack-channel-notification

## Install

require Golang environment and setup GOROOT.

```
$ go get github.com/hashibiroko/slack-channel-notification
```

## Usage

Please invite a bot in advance to the channel you want to notify.

#### Example 1:

```
$ slack-channel-notification -token=xxxxxx-xxxxxxxxx
```

#### Example 2: setting environment

```
$ export SLACK_BOT_TOKEN="xxxxxx-xxxxxxxxx"
$ slack-channel-notification
```

#### Example 3: using docker

```
$ docker run -itd --name slack-channel-notification -e SLACK_USER_TOKEN=xxxxxx-xxxxxxxxx hashibiroko/slack-channel-notification
```

### Flags

| name | description | default | require | environment |
| :--- | :---------- | :-----: | :-----: | :---------- |
| token | Set your slack bot token |  | true | SLACK_BOT_TOKEN |
| channel | Set your slack notification channel | random |  | SLACK_CHANNEL_NAME |
| delay | Set slack notification delay | 5 |  | DELAY |
