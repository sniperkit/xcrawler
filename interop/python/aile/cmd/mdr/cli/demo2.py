#!/usr/bin/env python2
from __future__ import absolute_import
from __future__ import print_function

import sys
import requests
from mdr import MDR
import mdr_extract
import re
from pathlib2 import Path

if __name__ == '__main__':
    if len(sys.argv) > 1:
        url = sys.argv[1]
    else:
        url = "http://www.yelp.co.uk/biz/the-ledbury-london"
    print("url: ", url)
    mdr = MDR()
    r = requests.get(url)
    doc = mdr.parseHtml(r.text.encode('utf8'))
    candidates = MDR().list_candidates(doc)
    for p in [doc.getpath(c) for c in candidates]: print(p)
    results = list(mdr_extract.extractResultsFromRoots(doc, candidates))

    r1 = "(US)\\d[,.']?\\d\\d\\d[,.']?\\d\\d\\d"
    r2 = "\\d[,.']\\d\\d\\d[,.']\\d\\d\\d"
    r = re.compile("(?<![^ :])((" + r1 + ")|(" + r2 + "))(?=[,.';:()]?(?![^ ]))")

    def textContainsPatent(text):
        return len(r.findall(text)) > 0

    results = [result for result in results if mdr_extract.someXmlTextInResultSatisfies(result, textContainsPatent)]

    for result in results:
        mdr_extract.printResult(result)

    for result in results:
        mdr_extract.printData(mdr_extract.extractDataFromResult(result))

    print("+++ END")