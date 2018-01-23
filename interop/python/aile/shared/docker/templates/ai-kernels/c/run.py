#! /usr/bin/env python3

import asyncio
import logging
import os
from pathlib import Path
import shlex
import sys
import tempfile

sys.path.insert(0, os.path.abspath('.'))
from base_run import BaseRunner
# For debugging
# sys.path.insert(0, os.path.abspath('..'))
# from base.run import BaseRun

log = logging.getLogger()

DEFAULT_CFLAGS = '-Wall'
DEFAULT_LDFLAGS = '-lrt -lm -pthread -ldl'
CHILD_ENV = {
    'TERM': 'xterm',
    'LANG': 'C.UTF-8',
    'SHELL': '/bin/ash',
    'USER': 'work',
    'HOME': '/home/work',
    'PATH': '/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin',
}


class CProgramRunner(BaseRunner):

    log_prefix = 'c-kernel'

    def __init__(self):
        super().__init__()
        self.child_env.update(CHILD_ENV)

    async def init_with_loop(self):
        self.user_input_queue = asyncio.Queue()

    async def build(self, build_cmd):
        if build_cmd is None or build_cmd == '':
            # skipped
            return
        elif build_cmd == '*':
            # use the default heuristic
            if Path('Makefile').is_file():
                await self.run_subproc('make')
            elif Path('main.c').is_file():
                cfiles = list(Path('.').glob('**/*.c'))
                ofiles = [Path(p.stem + '.o') for p in cfiles]
                for cf in cfiles:
                    cmd = f'gcc -c {cf} {DEFAULT_CFLAGS}'
                    await self.run_subproc(cmd)
                cfiles = ' '.join(map(lambda p: shlex.quote(str(p)), cfiles))
                ofiles = ' '.join(map(lambda p: shlex.quote(str(p)), ofiles))
                cmd = f'gcc {ofiles} {DEFAULT_LDFLAGS} -o ./main'
                await self.run_subproc(cmd)
            else:
                log.error('cannot find build script ("Makefile") '
                          'or the main file ("main.c").')
        else:
            await self.run_subproc(build_cmd)

    async def execute(self, exec_cmd):
        self.child_env.update({
            'LD_PRELOAD': os.environ.get('LD_PRELOAD', '/home/sorna/patch-libs.so'),
        })
        if exec_cmd is None or exec_cmd == '':
            # skipped
            return
        elif exec_cmd == '*':
            if Path('./main').is_file():
                await self.run_subproc('./main')
            elif Path('./a.out').is_file():
                await self.run_subproc('./a.out')
            else:
                log.error('cannot find executable ("a.out" or "main").')
        else:
            await self.run_subproc(exec_cmd)

    async def query(self, code_text):
        self.child_env.update({
            'LD_PRELOAD': os.environ.get('LD_PRELOAD', '/home/sorna/patch-libs.so'),
        })
        with tempfile.NamedTemporaryFile(suffix='.c', dir='.') as tmpf:
            tmpf.write(code_text.encode('utf8'))
            tmpf.flush()
            cmd = (f'gcc {tmpf.name} {DEFAULT_CFLAGS} -o ./main {DEFAULT_LDFLAGS} && '
                   f'./main')
            await self.run_subproc(cmd)

    async def complete(self, data):
        return []

    async def interrupt(self):
        # subproc interrupt is already handled by BaseRunner
        pass


if __name__ == '__main__':
    CProgramRunner().run()
