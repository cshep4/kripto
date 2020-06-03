import logging
import random
from decimal import Decimal

from model.rate import Rate


class Decider:
    def __init__(self, logger: logging.Logger):
        self.logger = logger

    def decide(self, rates: [Rate]) -> (bool, float, str):
        self.logger.info("deciding whether to trade...")

        # randomly decide whether to trade (1 in 60 chance, should average once an hour)
        decision = (random.randint(0,59) == 16)
        # amount is 1/100 BTC
        amount = round(Decimal(rates[0].rate / 100), 2)
        # randomly decide whether to buy or sell
        trade_type = random.choice(["buy", "sell"])

        return decision, amount, trade_type
