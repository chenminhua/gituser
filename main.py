#!/usr/bin/env python
import urllib2
import json
import sys
import math
import time
from operator import itemgetter
from prettytable import PrettyTable
from multiprocessing.dummy import Pool as ThreadPool

user_url = "https://api.github.com/users/{}"
repos_url = "https://api.github.com/users/{}/repos?page={}"


req_header = {'Connection': 'Keep-Alive',
    'Accept': 'text/html, application/xhtml+xml, */*',
    'Accept-Language': 'en-US,en;q=0.8,zh-Hans-CN;q=0.5,zh-Hans;q=0.3',
    'User-Agent': 'Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko'
 }


repos = []
retrieveDone = False

if len(sys.argv) < 2:
    print('Please input your github id!')
    sys.exit(1)
username = sys.argv[1]


print "retrieving {}'s info ...'".format(username)
user = json.load(urllib2.urlopen(urllib2.Request(user_url.format(username),None,req_header) ))
ut = PrettyTable(["name","email","location","follower","following", "created_at"])
ut.add_row([user["login"], user["email"], user["location"],user["followers"],user["following"],user["created_at"]])
print ut

total_workers = int(math.ceil(user["public_repos"] / 30.))

print "retrieving {}'s repos ...".format(username)

pages = [repos_url.format(username, page) for page in range(1,total_workers+1)]
pool = ThreadPool(total_workers)

drepos = pool.map(lambda url : json.load(urllib2.urlopen(urllib2.Request(url,None,req_header))), pages)
pool.close()
pool.join()
#print drepos

# flatten
repos = []
for r in drepos:
    repos += r

print "{} has {} repos!".format(username, len(repos))

repos_sorted_by_stars = sorted(repos, key=itemgetter('stargazers_count'), reverse=True)

pt = PrettyTable(["name","star","fork","html","description","language"])
pt.padding_width = 1

total_stars = 0
for repo in repos_sorted_by_stars:
    if repo['fork']:
        continue
    if repo['stargazers_count'] > 0:
        total_stars += repo['stargazers_count']
        pt.add_row([repo['name'],repo['stargazers_count'],repo['forks'],repo['html_url'], repo["description"], repo["language"]])

print pt

print "{} has got {} stars!".format(username, total_stars)
