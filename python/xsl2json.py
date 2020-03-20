import pandas as pd
import json
import sys

if __name__ == "__main__":
    iS = "news.xls"
    oS = "eksmo.json"
    if (len(sys.argv) == 3):
        iS = sys.argv[1]
        oS = sys.argv[2]
    listOfObjects = list()
    df = pd.read_excel(iS).T
    for i in df:
        tmp = {
            "Name" : df[i]['Название'],
            "Author" : df[i]['Автор'],
            "Cost" : df[i]['Розн.цена'],
            "Photo" : df[i]['Код'],
            "Publish" : df[i]['Издательство'],
            "Year" : df[i]['Год издания'],
            "FullName" : df[i]['Полное название'],
            "Sheets" : df[i]['Страниц'],
            "ISBN" : df[i]['ISBN']
        }
        listOfObjects.append(tmp)
    with open(oS, "w") as fout:
        json.dump(listOfObjects, fout)



