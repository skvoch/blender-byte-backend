import json
import sys
import copy
import re
from pymystem3 import Mystem

def seq(n): 
    for i in range(n): 
        yield "exits/exit{0}.json".format(i)


def res(n): 
    for i in range(n): 
        yield "results/result{0}.json".format(i)

if __name__ == "__main__":
    dataScience = None
    with open("exmo_books.json", encoding='utf-8') as readable:
        dataScience = (json.load(readable))
    m = Mystem()
    objects = None
    dictObject = dict()
    for exit in seq(41):
        with open(exit, encoding='utf-8') as readable:
            objects = (json.load(readable))
        count = 0
        for i in objects:
            dictObject[m.lemmatize(i["name"])[0]] = 1
        for i in objects:
            dictObject[m.lemmatize(i["name"])[0]] += 1
        for key, val in dictObject.items():
            print("{1}\t{0}".format(key, val))
        
    # print(count)
# 1393 / 873