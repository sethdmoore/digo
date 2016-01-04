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
  "consumer_key": "XXXX",
  "consumer_secret": "YYYY",
  "access_token": "ZZZZ",
  "access_token_secret": "123456",
  "stalking": {
    "discordapp": [
      "*"
    ],
    "nihilist_arbys": [
      "125728388009820160"
    ]
  }
}
```

field                | type   | description
---------------------|--------|---------
consumer_key         | string | Twitter Authentication
consumer_secret      | string | Twitter Authentication
access_token         | string | Twitter Authentication
access_token_secret  | string | Twitter Authentication
stalking             | object | Parent object for accounts to watch.
stalking.twitteruser | array  | user to watch for tweets, array is discord channels to post to.
