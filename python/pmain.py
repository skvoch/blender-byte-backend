from flask import Flask, json
import random
import logging

api = Flask(__name__)

@api.route('/random/', methods=['GET'])
def get_companies():
  response = {
    "RandomValue": str(random.randint(0,2024))
  }
  
  logging.info(response)
  return json.dumps(response)

if __name__ == '__main__':
    api.run(port=7100)