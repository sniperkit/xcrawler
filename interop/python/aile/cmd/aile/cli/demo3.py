#!/usr/bin/env python2

import time
import sys

import scrapely.htmlpage as hp

import aile.kernel
import aile.ptree
import aile.slybot_project

if __name__ == '__main__':
    if len(sys.argv) > 1:
        url = sys.argv[1]
    else:
        url = "https://news.ycombinator.com"

    if len(sys.argv) > 2:
        out_path = sys.argv[2]
    else:
        out_path = '/data/projects/slybot-project'
    print("out_path: ", out_path)

    print 'Downloading URL...',
    t1 = time.clock()
    page = hp.url_to_page(url)
    print 'done ({0}s)'.format(time.clock() - t1)

    print 'Extracting items...',
    t1 = time.clock()
    ie = aile.kernel.ItemExtract(aile.ptree.PageTree(page))
    print 'done ({0}s)'.format(time.clock() - t1)

    print 'Generating slybot project...',
    t1 = time.clock()
    aile.slybot_project.generate(ie, out_path)
    print 'done ({0}s)'.format(time.clock() - t1)
