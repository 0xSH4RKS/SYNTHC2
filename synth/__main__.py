#! /usr/bin/env python3

"""
Usage: SYN.py [-h] [-v] (teamserver) [<args>...]

options:
    -h, --help                   Show this help message and exit
    -v, --version                Show version
"""

from docopt import docopt
from synth import VERSION


def run():
    args = docopt(__doc__, version=VERSION, options_first=True)
    if args['teamserver']:
        import synth.core.ts.__main__ as teamserver
        teamserver.start(docopt(teamserver.__doc__, argv=args["<args>"]))
