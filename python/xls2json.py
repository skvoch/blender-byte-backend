import pandas as pd
import json
import sys

def main(iS, oS):
    listOfObjects = list()
    df = pd.read_excel(iS).T
    for i in df:
        tmp = {
            "Name" : df[i]['Название'],
            "Author" : df[i]['Автор'],
            "Cost" : df[i]['Розн.цена'],
            "Photo" : df[i]['Изображение'],
            "Publish" : df[i]['Издательство'],
            "Date" : str(df[i]['Дата пост.']),
            "FullName" : df[i]['Полное название'],
            "Sheets" : df[i]['Страниц'],
            "ISBN" : df[i]['ISBN'],
            "Topic" : df[i]["Тема"],
            "Code" : df[i]["Штрихкод"],
            "Series" : df[i]["Серия"],
            "Description" : df[i]["Описание"]
        }
        listOfObjects.append(tmp)
    with open(oS, "w", encoding='utf-8') as fout:
        json.dump(listOfObjects, fout)


def read(iS, Os):
    with open(iS, encoding='utf-8') as readable:
        parse_json = (json.load(readable))
    print(json.dumps(parse_json, indent=4, sort_keys=True, ensure_ascii=False))




if __name__ == "__main__":
    iS = "priceext.xls"
    oS = "eksmo.json"
    if (len(sys.argv) == 3):
        iS = sys.argv[1]
        oS = sys.argv[2]
    #main(iS, oS)
    read(oS, iS)

