#!/usr/bin/env python
import requests
import twitter
import os
import sys
import time
import json

DIGO_API = "http://127.0.0.1:8086/v1/message"
INTERVAL = 4600.0
CONFIG_FILE = os.path.join(os.path.dirname(__file__), "twitter-stalk.conf")

CONFIG_KEYS = ["consumer_key", "consumer_secret", "access_token", "stalking"]

def validate_config(config):
    """
    Ensure there are no missing configuration keys in the config
    """
    missing = []
    for key in CONFIG_KEYS:
        if key not in config:
            missing.append(key)
    if missing:
        print "Error! The following keys are missing from %s" % CONFIG_FILE
        print ", ".join(missing)
        sys.exit(2)


def fetch_config():
    """
    Open the config file and read it
    """
    with open(CONFIG_FILE, "r") as f:
        try:
            conf = json.load(f)
        except Exception as e:
            print "Could not load config file %s, %s" % (CONFIG_FILE, e)
            sys.exit(2)
    validate_config(conf)
    return conf


def configure_twitter(config):
    """
    Instantiate a twitter API object
    """
    try:
        t = twitter.Api(consumer_key=config["consumer_key"],
                        consumer_secret=config["consumer_secret"],
                        access_token_key=config["access_token"],
                        access_token_secret=config["access_token_secret"])
    except KeyError as e:
        print "Could  not configure twitter API: %s" % e
        print "Ensure the file %s exists, is proper JSON, and contains:"
        print '{"consumer_key": "AAA", "consumer_secret: "YYY",'
        print '"access_token": "YYY", "access_token_secret": "ZZZ",'
        print '"stalking": {"discordapp", ["*"]}'
        sys.exit(2)
    except Exception as e:
        print "Could not instantiate twitter API object: %s" % e
        sys.exit(2)

    return t


def stalk(t, account):
    statuses = []
    err = ""
    now = time.time()
    try:
        tweets = t.GetUserTimeline(screen_name=account,
                                     count=4,
                                     exclude_replies=True,
                                     include_rts=False)
    except Exception as e:
        print "Exception calling Twitter API. Probably rate limits"
        print "Could not utilize twitter API: %s" % e
        err = "Exception: %s" % e
        return [], err

    for status in tweets:
        tweet_time =  status.created_at_in_seconds
        delta = now - tweet_time
        if delta < INTERVAL:
            print u"New Tweet: @%s: %s" % (account, status.text)
            statuses.append(status)
        else:
            pass
            # print "skipped, too old: %s %s" % (delta, status.text)
    return statuses, err


def post_statuses(account, statuses, channels):
    h = {"content-type": "application/json"}

    for status in statuses:
        message = []
        id_str = str(status.id)
        src = "/".join(("https://twitter.com", account, "status", id_str))
        message.append("**@%s** - Twitter - %s" % (account, src))
        j = {"prefix": "", "payload": message, "channels": channels}
        try:
            r = requests.post(DIGO_API, headers=h, json=j)
            if r.status_code != 200:
                print "Could not reach the Digo API"
        except Exception as e:
            print "Exception contacting HTTP API: %s" % e
            continue


def main():
    """
    Fetch the config and the twitter API object
    Stalk the accounts specified forever...
    """
    backoff = INTERVAL
    errors = 0
    err = ""

    config = fetch_config()
    t = configure_twitter(config)

    print "Twitter Stalk is stalking..."
    while True:
        for account, channels in config["stalking"].iteritems():
            statuses, err = stalk(t, account)
            # exit loop immediately if rate limited
            if err:
                break
            post_statuses(account, statuses, channels)

        # exponential backoff
        # nice for tuning against rate limiting
        if err:
            errors += 1

            # ignore hiccups before tuning backoff
            if errors > 3:
                print "Tuning backoff from %s to %s" % (backoff, backoff * 2)
                backoff *= 2
                print "Sleeping for 15 minutes"
                print "https://dev.twitter.com/rest/public/rate-limiting"
                time.sleep(900)

            # don't want to sleep for 900 + backoff
            continue

        time.sleep(INTERVAL)


if __name__ == "__main__":
    main()
