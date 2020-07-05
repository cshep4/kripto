import logging
import boto3
import json

from decision.decider import Decider
from rate.retriever import Retriever
from wallet.retriever import Retriever as walletRetriever
from trade.trader import Trader

logger = logging.getLogger()
logger.setLevel(logging.INFO)

logger.info('initialisation')

lambda_client = boto3.client('lambda')

rateRetriever = Retriever(logger, lambda_client)
walletRetriever = walletRetriever(lambda_client)
decider = Decider(logger)
trader = Trader(logger, lambda_client)


def handler(event, context):
    logger.debug('Received event: {}'.format(event))

    rates = rateRetriever.get_rates()
    decision, amount, trade_type = decider.decide(rates)

    gbp, btc = walletRetriever.get_balances()

    if decision:
        trader.trade(amount, trade_type)
