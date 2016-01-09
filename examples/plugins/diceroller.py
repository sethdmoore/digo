#!/usr/bin/env python
import os
import sys
import json
import argparse
import random


config = \
{
    "name": "Dice roller",
    "triggers": ["/roll"],
    "type": "json",
    "description": "Roll diiiice!"
}


def setup_argparse():
    parser = argparse.ArgumentParser(prog='diceroller', add_help=False)
    parser.add_argument('cmd', choices=["register", "help", "json"], type=str, help='')
    parser.add_argument('json', nargs='?', type=str, help=None)
    return parser


def register():
    print json.dumps(config)


def determine_roll(dice):
    output = []
    if "d" in dice:
        num, size = dice.split("d")
        size = int(size)
        num = int(num)
        for dice in range(num):
            output.append(random.randint(1, size))
    else:
        return ""

    return ", ".join(map(str, output))


def parse_message(msg, response):
    output = ""
    successful_roll = False

    if "user" in msg:
        output += msg["user"] + " rolled "

    if "arguments" in msg:
        if len(msg["arguments"]) == 0 or len(msg["arguments"]) == 1 and \
                msg["arguments"][0] == "":
            output += determine_roll("1d6")
            response["payload"].append(output)
        else:
            for roll in msg["arguments"]:
                output += roll + ": "
                output += determine_roll(roll)
                response["payload"].append(output)
                output = ""



    if not response["payload"]:
        print_help(response)
        return


    print json.dumps(response)


def build_message(j):
    response = {}
    response["payload"] = []
    response["channels"] = []
    found_dice = False

    try:
        msg = json.loads(j)
    except Exception as e:
        response["payload"].append("Could not parse message from Digo: %s" % e)
        print json.dumps(response)
        sys.exit(2)

    if "channel" in msg:
        response["channels"].append(msg["channel"])

    if "help" in msg["arguments"]:
        print_help(msg, response)

    # input validatoin
    for arg in msg["arguments"]:
        if "d" in arg:
            if len(arg.split("d")) == 2:
                found_dice = True


    if found_dice:
        return msg, response
    else:
        print_help(msg, response)


def print_help(msg, response):
    response["payload"] = []
    response["payload"].append("==Dice Roller Help==")
    response["payload"].append("/roll  (rolls 1d6)")
    response["payload"].append("/roll 2d100  (rolls 2 d100's)")
    response["payload"].append("/roll 2d6 3d10 (rolls 2d6's and 3d10's)")
    print json.dumps(response)
    sys.exit(0)


def main():
    parser = setup_argparse()
    args = parser.parse_args()
    #print args
    if args.cmd == "register":
        register()
    elif args.cmd == "json" and args.json:
        msg, response = build_message(args.json)
        parse_message(msg, response)
    #print args

#        print_help()


if __name__ == '__main__':
    main()
