import json
import os


class Retriever:
    def __init__(self, client):
        self.client = client

    def get_rates(self) -> dict:
        resp = self.client.invoke(
            FunctionName=os.environ['READER_FUNCTION_NAME'],
            InvocationType='RequestResponse',
            Payload="")

        d = json.loads(resp['Payload'].read())

        rates: dict = {"p": []}

        for r in d:
            rates["p"].append(r['rate'])

        return rates
