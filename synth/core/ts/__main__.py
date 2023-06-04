#!/usr/bin/env python3

"""
Usage: teamserver [-h] [--port <PORT>] [--insecure] <host> <password>

optional arguments:
    -h, --help          Show this help message and exit
    -p, --port <PORT>   Port to bind to [default: 5000]
    --insecure          Start server without TLS
"""
import asyncio
from synth.core.ts.users import Users, UsernameAlreadyPresentError
from synth.core.ts.contexts import Listeners, Sessions, Modules, Stagers

class TeamServer:
    def __init__(self):
        self.users = Users()
        self.loop = asyncio.get_running_loop()
        self.contexts = {
            'listeners': Listeners(self),
            'sessions': Sessions(self),
            'modules': Modules(self),
            'stagers': Stagers(self),
            'users': self.users
        }

