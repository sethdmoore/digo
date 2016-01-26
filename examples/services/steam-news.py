#!/usr/bin/env python
import requests
from requests.auth import HTTPBasicAuth
import twitter
import os
import sys
import time
import json

DIGO_MSG_ROUTE = "v1/message"

ARTICLES = "3"
STEAM_API = \
        ("http://api.steampowered.com/ISteamNews/"
        "GetNewsForApp/v0002/?count=%s&format=json&appid=") % ARTICLES

INTERVAL = 360.0
CONFIG_FILE = os.path.join(os.path.dirname(__file__), "steam-news.json")

# example request
# http://api.steampowered.com/ISteamNews/GetNewsForApp/v0002/?appid=570&count=10&format=json

CONFIG_KEYS = ["steam_apps", "digo_api_url"]

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


def check_news(appid):
    statuses = []
    err = ""
    now = time.time()
    url = STEAM_API + appid
    #blob = {}
    "=".join(("appid", str(appid)))

    try:
        r = requests.get(STEAM_API + appid)

    except Exception as e:
        err = "Exception: %s" % e
        return [], err

    try:
        blob = r.json()

    except Exception as e:
        print "Did not receive JSON for appid: %s" % str(appid)
        return

    if not "appnews" in blob:
        print "No news??"
        print "%s" % blob
        return

    if not "newsitems" in blob["appnews"]:
        print "No newsitems??"
        print "%s" % blob
        return

    for item in blob["appnews"]["newsitems"]:
        # skip BS
        if item["is_external_url"]:
            continue

        dtime = item["date"]
        delta = now - dtime

        if delta < INTERVAL:
            print "New item: %s, %s" % (item["title"], item["url"])
            statuses.append(item)
        else:
            pass
            # print "skipped, too old: %s %s" % (delta, status.text)
    return statuses, err


def post_statuses(config, account, items, channels):
    h = {"content-type": "application/json"}
    auth_enabled = False
    api = "/".join((config["digo_api_url"], DIGO_MSG_ROUTE))


    if "basic_auth_user" in config and "basic_auth_password" in config:
        if config["basic_auth_password"] and config["basic_auth_user"]:
            auth_enabled = True

    if auth_enabled:
        user = config["basic_auth_user"]
        passwd = config["basic_auth_password"]
        auth = HTTPBasicAuth(user, passwd)

    for item in items:
        message = []
        # id_str = str(status.id)
        # src = "/".join(("https://twitter.com", account, "status", id_str))
        message.append("**%s** %s" % (item["title"], item["url"]))
        j = {"prefix": "", "payload": message, "channels": channels}
        try:
            if auth_enabled:
                r = requests.post(api, headers=h, json=j, auth=auth)
            else:
                r = requests.post(api, headers=h, json=j)

            if r.status_code == 200:
                print "Posted successfully to Digo API"
            elif r.status_code == 401:
                print "Received 401 Unauthorized from Digo API"
                print 'Please set "basic_auth_user" and "basic_auth_password"'
                print "in %s" % CONFIG_FILE
            else:
                print "Unhandled error hitting the Digo API"
                print "Received %s, expecting 200 or 401" % r.status_code
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

    print "Patrolling for steam news"
    while True:
        for app, channels in config["steam_apps"].iteritems():
            news, err = check_news(app)
            # exit loop immediately if rate limited
            if err:
                break
            post_statuses(config, app, news, channels)

        # exponential backoff
        # nice for tuning against rate limiting
        if err:
            errors += 1

            # ignore hiccups before tuning backoff
            if errors > 3:
                print "Tuning backoff from %s to %s" % (backoff, backoff * 2)
                backoff *= 2
                print "Sleeping for 15 minutes"
                time.sleep(900)

            # don't want to sleep for 900 + backoff
            continue

        time.sleep(INTERVAL)


if __name__ == "__main__":
    main()
