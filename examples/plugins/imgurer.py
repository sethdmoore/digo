#!/usr/bin/env python
import os
import sys
import json
import requests
import random

from imgurpython import ImgurClient

# DEAD SIMPLE PLUGIN

YOUTUBE_API_SERVICE_NAME = "youtube"
YOUTUBE_API_VERSION = "v3"

config = \
{
    "name": "imgurer",
    "triggers": ["/imgur", "/img"],
    "description": "Searches imgur for keywords"
}

def print_help():
    print "usage: /img keyword search"
    print "usage: /imgur random"


def no_results(args):
    print "No results for \"%s\"" % " ".join(args)
    sys.exit(0)



def fetch_conf():
    conf_file = os.path.join(os.path.dirname(__file__), "_imgurer.json")
    with open(conf_file, "r") as f:
        try:
            conf = json.load(f)
        except Exception as e:
            print "Could not read configuration file: %s" % e
            sys.exit(2)
        try:
            client_id = conf["client_id"]
            client_secret = conf["client_secret"]

            return client_id, client_secret
        except Exception as e:
            print "The config file is dead simple."
            print "{\"client_id\": \"XXXXX\", \"client_secret\": \"YYYY\"}."
            print "%s" % e
            sys.exit(2)

def filter_results(incoming_items, nsfw):
    items = []
    # filter out albums
    for item in incoming_items:
        addable = True
        if item.nsfw and nsfw:
            addable = True

        elif item.nsfw and not nsfw:
            addable = False

        if item.is_album:
            addable = False

        if addable:
            items.append(item)

    return items


def search(args, nsfw, client):
    term = "+".join(args)

    potential_items = client.gallery_search(term)
    items = filter_results(potential_items, nsfw)
    try:
        item = random.choice(items)
    except Exception as e:
        no_results(args)
    print item.link


def rando(nsfw, client):
    potential_items = client.gallery_random()
    items = filter_results(potential_items, nsfw)
    try:
        item = random.choice(items)
    except Exception as e:
        print "IMGUR API problem? No random image :[]"
        sys.exit(0)
    print item.link


def main():
    client_id, client_secret = fetch_conf()
    client = ImgurClient(client_id, client_secret)
    if len(sys.argv) == 2:
        if sys.argv[1] == "register":
            print json.dumps(config)
        elif sys.argv[1] == "random":
            rando(False, client)
        elif sys.argv[1] == "help":
            print_help()
        elif sys.argv[1] == "randomnsfw":
            rando(True, client)
        else:
            search(sys.argv[1:], False, client)
    elif len(sys.argv) > 2:
        if sys.argv[1] == "searchnsfw":
            search(sys.argv[2:], True, client)
        elif sys.argv[1] == "search":
            search(sys.argv[2:], False, client)
        else:
            search(sys.argv[1:], False, client)
    elif len(sys.argv) > 1:
        search(sys.argv[1:], False, client)
    else:
        print_help()


if __name__ == "__main__":
    main()
