import logging
import json
import os


class Retriever:
    def __init__(self, logger: logging.Logger, client):
        self.logger = logger
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

        self.logger.info('{}'.format(rates))

        return rates
