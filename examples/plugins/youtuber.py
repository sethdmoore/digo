#!/usr/bin/env python
import os
import sys
import json
import requests

from googleapiclient.discovery import build
from oauth2client.tools import argparser

# DEAD SIMPLE PLUGIN
# call it with

YOUTUBE_API_SERVICE_NAME = "youtube"
YOUTUBE_API_VERSION = "v3"

config = \
{
    "name": "youtuber",
    "triggers": ["/youtube", "/yt"],
    "description": "Searches youtube for keywords"
}

def print_help():
    print "usage: /yt keyword search"


def fetch_key():
    conf_file = os.path.join(os.path.dirname(__file__), "_youtuber.conf")
    with open(conf_file, "r") as f:
        try:
            conf = json.load(f)
        except Exception as e:
            print "Could not read configuration file: %s" % e
            sys.exit(2)
        try:
            api_key = conf["api_key"]
            return api_key
        except Exception as e:
            print "The config file is dead simple: {\"api_key\": \"XXXXX\"}"
            sys.exit(2)

def search(args):
    api_key = fetch_key()
    youtube = build(YOUTUBE_API_SERVICE_NAME, YOUTUBE_API_VERSION, developerKey=api_key)
    search_response = youtube.search().list(
            q="+".join(args),
            part="id,snippet",
            maxResults=10
            ).execute()
    results = search_response.get("items", [])
    for result in results:
        vid = result["id"]
        if vid["kind"] == "youtube#video":
            print "https://youtu.be/" + vid["videoId"]
            sys.exit(0)
    else:
        print "No results found for query \"%s\"" % " ".join(args) 
        sys.exit(0)

def main():
    try:
        if sys.argv[1] == "register":
            print json.dumps(config)
        elif sys.argv[1] == "help":
            print_help
        else:
            search(sys.argv[1:])
    except IndexError as e:
        print_help()


if __name__ == "__main__":
    main()
