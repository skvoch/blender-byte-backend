import pandas as pd
import json
import sys

def main(iS, oS):
    listOfObjects = list()
    df = pd.read_excel(iS).T
    df = df.fillna("")
    for i in df:
        tmp = {
            "Name" : df[i]['Название'],
            "Author" : df[i]['Автор'],
            "Cost" : df[i]['Розн.цена'] and 0,
            "Photo" : df[i]['Изображение'].replace("\\", "%5C"),
            "Publish" : df[i]['Издательство'],
            "Date" : str(df[i]['Дата пост.']),
            "FullName" : df[i]['Полное название'],
            "Sheets" : df[i]['Страниц'] and 0,
            "ISBN" : df[i]['ISBN'],
            "Topic" : df[i]["Тема"],
            "Code" : df[i]["Штрихкод"] and 0,
            "Series" : df[i]["Серия"] and 0,
            "Description" : df[i]["Описание"]
        }
        listOfObjects.append(tmp)
    with open(oS, "w", encoding='utf-8') as fout:
        json.dump(listOfObjects, fout, ensure_ascii=False, allow_nan=False )


def read(iS, Os):
    with open(iS, encoding='utf-8') as readable:
        parse_json = (json.load(readable))
    print(parse_json[28])




if __name__ == "__main__":
    iS = "priceext.xls"
    oS = "exmo_books.json"
    if (len(sys.argv) == 3):
        iS = sys.argv[1]
        oS = sys.argv[2]
    #main(iS, oS)
    read(oS, iS)

