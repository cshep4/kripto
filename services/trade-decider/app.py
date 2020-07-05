import logging
# import os
import boto3
import json

# from lumigo_tracer import lumigo_tracer
from decision.decider import Decider
from rate.retriever import Retriever
from wallet.retriever import Retriever as walletRetriever
from trade.trader import Trader

# Set up logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)

logger.info('initialisation')

# dynamo_client = boto3.client('dynamodb')
lambda_client = boto3.client('lambda')

rateRetriever = Retriever(logger, lambda_client)
walletRetriever = walletRetriever(lambda_client)
decider = Decider(logger)
trader = Trader(logger, lambda_client)


# @lumigo_tracer(token=os.environ['LUMIGO_TRACER_TOKEN'], enhance_print=True)
def handler(event, context):
    logger.debug('Received event: {}'.format(event))

    rates = rateRetriever.get_rates()
    decision, amount, trade_type = decider.decide(rates)

    gbp, btc = walletRetriever.get_balances()

    if decision:
        trader.trade(amount, trade_type)
    # for record in event['Records']:
    #     payload = json.loads(record['body'], parse_float=str)
    #     operation = record['messageAttributes']['Method']['stringValue']
    #     if operation in operations:
    #         try:
    #             operations[operation](dynamo_client, payload)
    #             logger.info('{} method successful'.format(operation))
    #             logger.debug('Event payload: {}'.format(payload))
    #         except Exception as e:
    #             logger.error(e)
    #     else:
    #         logger.error('Unsupported method \'{}\''.format(operation))
