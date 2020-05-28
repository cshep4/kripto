import logging
import json
import os


class Trader:
    def __init__(self, logger: logging.Logger, client):
        self.logger = logger
        self.client = client

    def trade(self, amount: float, trade_type: str):
        self.logger.info(trade_type + "ing " + str(amount))

        self.client.invoke(
            FunctionName=os.environ['TRADER_FUNCTION_NAME'],
            InvocationType='RequestResponse',
            Payload=json.dumps({"tradeType": trade_type, "amount": str(amount)})
        )
