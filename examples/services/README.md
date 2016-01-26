# twitter-stalk.py
A simple daemon that follows users and their tweets.
Follow any number (within rate limit reason) of twitter users and post their tweets to whichever channel (individually configurable).

## Features
* exponential rate limit backoff
* configured through JSON
* configure which user's tweets go to which channel(s)

## Configuration
A json file named "twitter-stalk.conf" is expected to live in the same directory as the service.  
Since it's written in Python and this is only an "example", this is easily changed in the code

```json
{
  "digo_api_url": "http://127.0.0.1:8086",
  "basic_auth_user": "mydigobasicauth",
  "basic_auth_password": "mydigopassword",
  "consumer_key": "XXXX",
  "consumer_secret": "YYYY",
  "access_token": "ZZZZ",
  "access_token_secret": "123456",
  "stalking": {
    "discordapp": [
      "*"
    ],
    "nihilist_arbys": [
      "123456789012345678"
    ]
  }
}
```

field                | type   | description
---------------------|--------|---------
digo_api_url         | string | Digo API URL and port
basic_auth_user      | string | Digo Basic Auth, optional
basic_auth_password  | string | Digo Basic Auth password, optional
consumer_key         | string | Twitter Authentication
consumer_secret      | string | Twitter Authentication
access_token         | string | Twitter Authentication
access_token_secret  | string | Twitter Authentication
stalking             | object | Parent object for accounts to watch.
stalking.twitteruser | string | user to watch for tweets, array is discord channels to post to.



# steam-news.py
A simple daemon that follows the news page of a steam application

## Features
* exponential rate limit backoff
* configured through JSON
* configure which application news go to which channel

## Configuration
A json file named "steam-news.json" is expected to live in the same directory as the service.  
Since it's written in Python and this is only an "example", this is easily changed in the code

```json
{
  "digo_api_url": "http://127.0.0.1:8086",
  "basic_auth_user": "mydigobasicauth",
  "basic_auth_password": "mydigopassword",
  "steam_apps": {
    "570": [
      "2345676829"
    ]
  }
}
```

field                | type   | description
---------------------|--------|---------
digo_api_url         | string | Digo API URL and port
basic_auth_user      | string | Digo Basic Auth, optional
basic_auth_password  | string | Digo Basic Auth password, optional
steam_apps           | object | Parent object for items to watch.
steam_apps.id        | string | Key is the steamapp id, value is the channels to post to
