#!/usr/bin/env python
import requests
import json
import sys

config = \
{
    "name": "Urban Dictionary",
    "triggers": ["/urban", "/ud", "/urbandictionary"],
    "description": "Searches for Urban Dictionary terms"
}


def register():
    print json.dumps(config)


def print_help():
    print "Usage: %s search terms" % config["triggers"][0]


def query(q):
    terms = "+".join(q)
    try:
        r = requests.get("http://api.urbandictionary.com/v0/define?term=" + terms)
    except Exception as e:
        print "Could not reach Urban Dictionary API"

    if not r.status_code == 200:
        print "General API failure"

    j = r.json()

    if j["result_type"] == "no_results":
        print "No results for \"%s\"" % " ".join(q)
        return
    elif j["result_type"] == "exact":
        word = j["list"][0]
        #print "Urban Dictionary"
        print ":  **%s** - Urban Dictionary" % word["word"]
        print "```%s```" % word["definition"]


def main():
    if not len(sys.argv) > 1:
        print_help()
    elif len(sys.argv) == 2 and sys.argv[1] == "register":
        register()
    else:
        args = sys.argv[1:]
        query(args)


if __name__ == "__main__":
    main()
