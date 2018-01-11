#!/usr/bin/env python2

import sys
import requests
from mdr import MDR

if __name__ == '__main__':
    if len(sys.argv) > 1:
        url = sys.argv[1]
    else:
        url = "http://www.yelp.co.uk/biz/the-ledbury-london"
    print("url: ", url)
    mdr = MDR()
    r = requests.get(url)
    candidates, doc = mdr.list_candidates(r.text.encode('utf8'))
    for c in candidates[:20]:
    	print doc.getpath(c)
    # print [doc.getpath(c) for c in candidates[:10]]