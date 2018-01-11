#!/usr/bin/env python2

import time
import sys
import codecs
import cgi

import scrapely.htmlpage as hp
import numpy as np

import aile.kernel
import aile.ptree


def annotate(page, labels, out_path="/data/annotated/annotated.html"):
    match = aile.ptree.match_fragments(page.parsed_body)
    with codecs.open(out_path, 'w', encoding='utf-8') as out:
        indent = 0
        for i, (fragment, label) in enumerate(
                zip(page.parsed_body, labels)):
            #if label >= 0:
            #    print('<span class="line" style="color:red">')
            #else:
            #    print('<span class="line" style="color:black">')
            if isinstance(fragment, hp.HtmlTag):
                if fragment.tag_type == hp.HtmlTagType.CLOSE_TAG:
                    if match[i] >= 0 and indent > 0:
                        indent -= 1
                    # print('{0:3d}|{1}'.format(label, fragment.tag))
                else:
                    # print('{0:3d}|{1}'.format(label, fragment.tag))
                    for k,v in fragment.attributes.iteritems():
                        print(u' {0}="{1}"'.format(k, v))
                    if fragment.tag_type == hp.HtmlTagType.UNPAIRED_TAG:
                        print('/')
                    if match[i] >= 0:
                        indent += 1
            else:
                print(u'{0:3d}|{1}'.format(label, cgi.escape(page.body[fragment.start:fragment.end].strip())))

if __name__ == '__main__':
    if len(sys.argv) > 1:
	    url = sys.argv[1]
    else:
    	url = "https://news.ycombinator.com"
    print("url: ", url)

    print 'Downloading URL...',
    t1 = time.clock()
    page = hp.url_to_page(url)
    print 'done ({0}s)'.format(time.clock() - t1)

    print 'Extracting items...',
    t1 = time.clock()
    ie = aile.kernel.ItemExtract(aile.ptree.PageTree(page))
    print 'done ({0}s)'.format(time.clock() - t1)

    print 'Annotating HTML'
    labels = np.repeat(-1, len(ie.page_tree.page.parsed_body))
    items, cells = ie.table_fragments[0]
    for i in range(cells.shape[0]):
        for j in range(cells.shape[1]):
            labels[cells[i, j]] = j
    annotate(ie.page_tree.page, labels)
