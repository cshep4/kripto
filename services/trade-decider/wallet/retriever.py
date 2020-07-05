import json
import os


class Retriever:
    def __init__(self, client):
        self.client = client

    def get_balances(self) -> (float, float):
        resp = self.client.invoke(
            FunctionName=os.environ['GET_WALLET_FUNCTION_NAME'],
            InvocationType='RequestResponse',
            Payload="")

        d = json.loads(resp['Payload'].read())

        return d["gbp"]["available"], d["btc"]["available"]
