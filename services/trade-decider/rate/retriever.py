import logging
import json
import os

from model.rate import Rate


class Retriever:
    def __init__(self, logger: logging.Logger, client):
        self.logger = logger
        self.client = client

    def get_rates(self) -> [Rate]:
        resp = self.client.invoke(
            FunctionName=os.environ['READER_FUNCTION_NAME'],
            InvocationType='RequestResponse',
            Payload="")

        d = json.loads(resp['Payload'].read())
        self.logger.info('{}'.format(d))

        return d
