#!/usr/bin/env python
import os
import sys
import json
__version__ = '0.3'
__releasedate__ = '2009-12-23'
__author__ = 'Ryan McGreal <ryan@quandyfactory.com>'
__homepage__ = 'http://quandyfactory.com/insult'
__repository__ = 'http://gist.github.com/258915'
__copyright__ = 'Copyright (C) 2009 by Ryan McGreal. Licenced under GPL version 2. http://www.gnu.org/licenses/gpl-2.0.html'


config = \
{
    "name": "Insult Generator",
    "triggers": ["/insult"],
    "description": "Generate an insult, directed at someone."
}


def register():
    print json.dumps(config)

def print_help():
    print "Usage: %s person" % config["triggers"][0]
    sys.exit(0)


def generate_insult(target):
    """
    Generates an Elizabethan insult.
    Insult terms are via: http://www.museangel.net/insult.html#generator
    Get insulted on the web at: http://quandyfactory.com/insult
    Update: Duh! The Elizabethan second person singular pronoun is "thou", not "you"
    """
    from random import randint
    words = (
        ('Artless', 'Bawdy', 'Beslubbering', 'Bootless', 'Churlish', 'Cockered', 'Clouted', 'Craven', 'Currish', 'Dankish', 'Dissembling', 'Droning', 'Errant', 'Fawning', 'Fobbing', 'Froward', 'Frothy', 'Gleeking', 'Goatish', 'Gorbellied', 'Impertinent', 'Infectious', 'Jarring', 'Loggerheaded', 'Lumpish', 'Mammering', 'Mangled', 'Mewling', 'Paunchy', 'Pribbling', 'Puking', 'Puny', 'Quailing', 'Rank', 'Reeky', 'Roguish', 'Ruttish', 'Saucy', 'Spleeny', 'Spongy', 'Surly', 'Tottering', 'Unmuzzled', 'Vain', 'Venomed', 'Villainous', 'Warped', 'Wayward', 'Weedy', 'Yeasty',),
        ('Base-court', 'Bat-fowling', 'Beef-witted', 'Beetle-headed', 'Boil-brained', 'Clapper-clawed', 'Clay-brained', 'Common-kissing', 'Crook-pated', 'Dismal-dreaming', 'Dizzy-eyed', 'Dog-hearted', 'Dread-bolted', 'Earth-vexing', 'Elf-skinned', 'Fat-kidneyed', 'Fen-sucked', 'Flap-mouthed', 'Fly-bitten', 'Folly-fallen', 'Fool-born', 'Full-gorged', 'Guts-griping', 'Half-faced', 'Hasty-witted', 'Hedge-born', 'Hell-hated', 'Idle-headed', 'Ill-breeding', 'Ill-nurtured', 'Knotty-pated', 'Milk-livered', 'Motley-minded', 'Onion-eyed', 'Plume-plucked', 'Pottle-deep', 'Pox-marked', 'Reeling-ripe', 'Rough-hewn','Rude-growing', 'Rump-fed', 'Shard-borne', 'Sheep-biting', 'Spur-galled', 'Swag-bellied', 'Tardy-gaited', 'Tickle-brained', 'Toad-spotted', 'Unchin-snouted', 'Weather-bitten',),
        ('Apple-john', 'Baggage', 'Barnacle', 'Bladder', 'Boar-pig', 'Bugbear', 'Bum-bailey', 'Canker-blossom', 'Clack-dish', 'Clot-pole', 'Coxcomb', 'Codpiece', 'Death-token', 'Dewberry', 'Flap-dragon', 'Flax-wench', 'Flirt-gill', 'Foot-licker', 'Fustilarian', 'Giglet', 'Gudgeon', 'Haggard', 'Harpy', 'Hedge-pig', 'Horn-beast', 'Huggermugger', 'Jolt-head', 'Lewdster', 'Lout', 'Maggot-pie', 'Malt-worm', 'Mammet', 'Measle', 'Minnow','Miscreant', 'Mold-warp', 'Mumble-news', 'Nut-hook', 'Pigeon-egg', 'Pignut', 'Puttock','Pumpion', 'Rats-bane', 'Scut', 'Skains-mate', 'Strumpet', 'Varlot', 'Vassal', 'Whey-face', 'Wagtail',),
        )
    insult_list = (
        words[0][randint(0,len(words[0])-1)],
        words[1][randint(0,len(words[1])-1)],
        words[2][randint(0,len(words[2])-1)],
        )
    vowels = 'AEIOU'
    article = 'an' if insult_list[0][0] in vowels else 'a'
    return '%s, thou art %s %s, %s %s.' % (target, article, insult_list[0], insult_list[1], insult_list[2])


def main():
    if len(sys.argv) == 2 and sys.argv[1] == "register":
        register()
    elif len(sys.argv) == 2:
        print generate_insult(sys.argv[1])
    else:
        print_help()


if __name__ == '__main__':
    main()
