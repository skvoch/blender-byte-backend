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
    oS = "exmo_books.json"
    os = "exit.json"
    coff = None
    with open(oS, encoding='utf-8') as readable:
        coff = (json.load(readable))
    datascience = None
    with open(os, encoding='utf-8') as readable:
        datascience = (json.load(readable))
    result = list()
    m = Mystem()
        
    cur = 0
    tmp = None
    for i in datascience:
        try:
            if cur != int(i["id"]):
                tmp = copy.copy(coff[int(i["id"])])
                tmp["dataScience"] = [m.lemmatize(i["name"])[0], ]
            else:
                tmp["dataScience"].append(m.lemmatize(i["name"])[0])
            result.append(tmp)
        except TypeError:
            print("+1")
    print(json.dumps(result, ensure_ascii=False, indent=4))

    with open("result.json", "w", encoding='utf-8') as write:
        json.dump(result, write, ensure_ascii=False, allow_nan=False )

