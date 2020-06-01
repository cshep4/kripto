import logging
import random

from model.rate import Rate


class Decider:
    def __init__(self, logger: logging.Logger):
        self.logger = logger

    def decide(self, rates: [Rate]) -> (bool, float, str):
        self.logger.info("deciding whether to trade...")

        # randomly decide whether to trade
        decision = random.choice([True, False])
        # amount is 1/100 BTC
        amount = round(rates[0].rate / 100, 2)
        # randomly decide whether to buy or sell
        trade_type = random.choice(["buy", "sell"])

        return decision, amount, trade_type
