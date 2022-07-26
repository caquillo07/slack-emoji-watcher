# slack-emoji-watcher

## Slack setup instructions
 - [Create a new Slack app](https://api.slack.com/apps?new_app=1) and select the desired workplace
 - Here you can update the bot's `Display Information` before or after configuring.
 - Under `Add features and functionality`, we will give the bot the following:
   - Event Subscriptions (required for listening to `emoji_changed` events)
   - Bot
   - Permissions (required for sending messages to a channel)

### Event Subscription and app-level token
 - Go to `Settings` -> `Socket Mode` and click the enable toggle.
     - This will generate an `app-level` token, save this as it will be required later
 - Under the `Features Affected` table, click `Event Subscriptions` (can also from the left pane `Features` -> `Event Subscription`)
 - Toggle `Enable Events` to `On`
 - Under `Subscribe to bot events`, add `emoji_changed`
 - Click `Save Changes` at the bottom of the screen, and return to `Settings` -> `Basic Information`

### Bot
In this page (also reached by going to `Features` -> `App Home`) you are able to change the bot's basic information if needed

### Permissions
 - Go to `Features` -> `OAuth & Permissions`
 - Under the `Scopes` section, give the `chat:write` scope under `Bot Token Scopes`, this will the bot send messages to the given channel
 - Finally, at the top of this page click `Install to Work` under `OAuth Tokens for Your Workspace` and follow the prompts.
   - This step will generate the `Bot User OAuth Token`, save this as it will also be required later.

### Using the Bot
After setting up, the bot should be displayed under the `Apps` section in slack. From here, invite the bot to the desired channel.

After running the bot, it should connect and start posting to slack when a new emoji is added.


## Configuration
The bot can be configured by setting the follow environment variables:
 - ENV: `env|prod` defaults to `env`
   - Setting this value to `prod` will produce logs in JSON form, leaving as `dev` will produce pretty logs
 - SLACK_APP_TOKEN: the `xapp` token generated above, this field is required
 - SLACK_BOT_TOKEN: the `xoxb` token generated above, this field is required
 - EMOJI_CHANNEL: the channel in which the bot will post, this defaults to `#general`
   - This field can be added with and without a `#` (i.e. `#emoji` and `emoji` are both valid entries)

## Building and running the bot
```shell
go build -o emoji-bot *.go

# if building for linux x86
GOOS=linux GOARCH=amd64 go build -o emoji-bot *.go

# run the bot
./emoji-bot
```

once running, the bot will receive a hello message from Slack
```shell
5:59PM INF Incoming WebSocket message: {
  "type": "hello",
  "num_connections": 1,
  "debug_info": {
    "host": "applink-<some-hash>",
    "build_number": 9,
    "approximate_connection_time": 18060
  },
  "connection_info": {
    "app_id": "abc1234"
  }
}`
```
