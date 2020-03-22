# Imports the Google Cloud client library
from google.cloud import language
from google.cloud.language import enums
from google.cloud.language import types
import pandas as pd
import json
import sys
import copy
import time
from google.api_core.exceptions import InternalServerError
client = language.LanguageServiceClient()

def getCurNot(descript, df):
    my_types = ["UNKNOWN", "PERSON", "LOCATION", "ORGANIZATION", "EVENT",
        "WORK_OF_ART", "CONSUMER_GOOD", "OTHER", "PHONE_NUMBER", "ADDRESS", "DATE", 
        "NUMBER", "PRICE"]

    global client
    try:
        document = types.Document(
            content=descript,
            language="ru",
            type=enums.Document.Type.PLAIN_TEXT
        )
    except TypeError:
        print("error")
        return {}
    try:
        response = client.analyze_entities(
            document=document,
            encoding_type="UTF32"
        )
    except InternalServerError:
        print("error")
        return {}
    listOfNotify = list()
    for entity in response.entities:
        listOfNotify.append(entity.name)
        if(entity.salience < 0.05):
            break
    df["DataScience"] = listOfNotify
    return df
def addNotation(df):
    asd = df.T
    for i in asd:
        (yield getCurNot(asd[i]["Description"], asd[i])) if i else (yield {}) 
def seq(n): 
    for i in range(0, n): 
        yield "exits/exit{0}.json".format(i)
def start():
    iS = "exmo_books.json"
    oS = [i for i in seq(41)]
    parse_json = None
    with open(iS, encoding='utf-8') as readable:
        parse_json = (json.load(readable))
    df = pd.DataFrame.from_dict(parse_json)
    count = 0
    for dataflow in oS:
        count += 1
        tmp = []
        print("{0}/{1}".format(count, 41))
        start = time.perf_counter()
        for i in addNotation(df[count * 1000:
            (count + 1) * 1000 if (count + 1) * 1000 < len(df) else -1]) :
                if type(i) != dict:
                    #print(i)
                    tmp.append(i)
        with open(dataflow, "w", encoding='utf-8') as write:
            pd.DataFrame(tmp).to_json(dataflow, force_ascii = False,orient='records')
            # json.dump(tmp, write, ensure_ascii=False, allow_nan=False )
        stop = time.perf_counter()
        print("time :",  stop - start)
start()
