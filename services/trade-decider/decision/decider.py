import logging
import random
from decimal import Decimal

from model.rate import Rate


class Decider:
    def __init__(self, logger: logging.Logger):
        self.logger = logger
        self.rolling_window = 100
        self.base_column = "p" # price/rate
        self.ema_short_span = 9
        self.ema_long_span = 21
        
        
    def decide(self, rates: [Rate]) -> (bool, float, str):
        self.logger.info("deciding whether to trade...")

        # randomly decide whether to trade (1 in 60 chance, should average once an hour)
        decision = (random.randint(0,59) == 16)
        # amount is 1/100 BTC
        amount = round(Decimal(rates[0].rate / 100), 2)
        # randomly decide whether to buy or sell
        trade_type = random.choice(["buy", "sell"])

        return decision, amount, trade_type

    def decide_new(
        data: dict,
        bitcoin,
        usd,
        trade_percentage,
        trade_commission
            )-> (bool, float, str, dict):
        decision = False
        df_limited = pd.DataFrame.from_dict(data)
        df_limited['rolling'] = df_limited[self.base_column].rolling(self.rolling_window).mean()
        df_limited['ema_short'] = df_limited['rolling'].ewm(span=self.ema_short_span).mean()
        df_limited['ema_long'] = df_limited['rolling'].ewm(span=self.ema_long_span).mean()

        row = df_limited.iloc[-2]
        if row['ema_short']>row['ema_long']:
            ontop_prev = "s"
        else:
            ontop_prev = "l"
        ema_short_prev = row["ema_short"]
        ema_long_prev = row["ema_long"]
        # compare to now
        row = df_limited.iloc[-1]
        price = row["p"]
        trade_type = None
        amount = 0

        if row['ema_short']>row['ema_long']:
            ontop_now = 's'
        else:
            ontop_now = 'l'
        ema_short_now = row["ema_short"]
        ema_long_now = row["ema_long"]

        if ontop_now == 's' and ontop_prev == 'l':
            trade_type = "buy"
            decision = True
            amount = buy_bitcoin(usd,trade_percentage,price)

        if ontop_now == 'l' and ontop_prev == 's':
            decision = True
            trade_type = "sell"
            amount = buy_USD(bitcoin,trade_percentage,price)
        properties = {
            "ontop_now":ontop_now,
            "ontop_prev":ontop_prev,
            "price":price,
            "ema_short_prev":ema_short_prev,
            "ema_long_prev":ema_long_prev

        }
        return decision, amount, trade_type, properties
    
# def buy_USD(bitcoin,trade_percentage,price,trade_commission):
#     how_many = ((bitcoin*trade_percentage)*price)*(1-trade_commission)
#     return how_many

# def buy_bitcoin(usd,trade_percentage,price,trade_commission):
#     how_many = ((usd*trade_percentage)/price)*(1-trade_commission)
#     return how_many

def buy_USD(bitcoin,trade_percentage,price):
    how_many = ((bitcoin*trade_percentage)*price)
    return how_many

def buy_bitcoin(usd,trade_percentage,price):
    how_many = ((usd*trade_percentage)/price)
    return how_many