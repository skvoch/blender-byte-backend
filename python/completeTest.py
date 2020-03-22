import pandas as pd
import json
import sys
import copy
from pymystem3 import Mystem

def read(iS, Os):
    with open(iS, encoding='utf-8') as readable:
        parse_json = (json.load(readable))
    print(json.dumps(parse_json, indent=4, sort_keys=True))

if __name__ == "__main__":
    oS = "results1.json"
    coff = None
    with open(oS,encoding='utf-8') as readable:
        coff = json.load(readable)
    m = Mystem()
    tmpList = []
    tmp = {}
    for i in coff:
        for k in i["DataScience"]:
            k = m.lemmatize(k.lower())[0]
            tmp[k] = 0
            tmpList.append(k)
    for i in tmpList:
        tmp[i] += 1
    for key, val in tmp.items():
        print("{0} : {1}".format(key, val))

    with open("results3.json", "w", encoding='utf-8') as write:
        json.dump(coff, write, ensure_ascii=False, allow_nan=False )

