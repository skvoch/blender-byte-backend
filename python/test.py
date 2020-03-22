# Imports the Google Cloud client library
from google.cloud import language
from google.cloud.language import enums
from google.cloud.language import types
import pandas as pd
import json
import sys
# Instantiates a client
client = language.LanguageServiceClient()

# The text to analyze
text = u'''Берём грустную Наташку, Спрашиваем что случилось, Если не отвечает, можно чуть-чуть надавить на доверие : \"солнышко, ты же знаешь, 
что ты всегда можешь все мне рассказать. Я тебя выслушаю и постараюсь помочь\"
Два варианта 
1. После этой фразы сдаётся и рассказывает (потому что это очень мило и заботливо) 
2. Не рассказывает (просто вредничает). Можно ласково сказать \" вреднюлька моя, не хочешь не рассказывай 😛\".  Все равно сдаться и расскажет
'''
document = types.Document(
    content=text,
    language="ru",
    type=enums.Document.Type.PLAIN_TEXT)
response = client.analyze_entities(
    document=document,
    encoding_type="UTF32"
)

# Detects the sentiment of the text
#sentiment = client.analyze_sentiment(document=document).document_sentimentse
for entity in response.entities:
    print("++" * 20)
    print("name {0}".format(entity.name))
    print("type {0}".format(entity.type))
    print("metadata {0}".format(entity.metadata))
    print("salience {0}".format(entity.salience))



#print('Text: {}'.format(text))
#print('Sentiment: {}, {}'.format(sentiment.score, sentiment.magnitude))



# def main(iS, oS):
#     listOfObjects = list()
#     df = pd.read_excel(iS).T
#     df = df.fillna("")
#     for i in df:
#         tmp = {
#             "Name" : df[i]['Название'],
#             "Author" : df[i]['Автор'],
#             "Cost" : df[i]['Розн.цена'] and 0,
#             "Photo" : df[i]['Изображение'],
#             "Publish" : df[i]['Издательство'],
#             "Date" : str(df[i]['Дата пост.']),
#             "FullName" : df[i]['Полное название'],
#             "Sheets" : df[i]['Страниц'] and 0,
#             "ISBN" : df[i]['ISBN'],
#             "Topic" : df[i]["Тема"],
#             "Code" : df[i]["Штрихкод"] and 0,
#             "Series" : df[i]["Серия"] and 0,
#             "Description" : df[i]["Описание"]
#         }
#         listOfObjects.append(tmp)
#     with open(oS, "w", encoding='utf-8') as fout:
#         json.dump(listOfObjects, fout, ensure_ascii=False, allow_nan=False )

def addNotation(df):
    for i in df.T:
        



if __name__ == "__main__":
    iS = "priceext.xls"
    oS = "exmo_books.json"
    if (len(sys.argv) == 3):
        iS = sys.argv[1]
        oS = sys.argv[2]
    listOfObjects = list()
    df = pd.read_excel(iS)
    df = df.fillna(a)
    df = df + addNotation(df)