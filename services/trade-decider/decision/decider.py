import logging
import pandas as pd


class Decider:
    def __init__(self, logger: logging.Logger):
        self.logger = logger
        self.rolling_window = 60
        self.base_column = "rate" # price/rate
        self.ema_short_span = 9*60*24
        self.ema_long_span = 21*60*24

    def decide(self,
        data: dict,
        bitcoin,
        usd,
        trade_percentage,
        trade_commission
            )-> (bool, float, str, dict):
        decision = False
        df_limited = pd.DataFrame.from_dict(data)
        df_limited['rolling'] = df_limited[self.base_column].rolling(self.rolling_window).mean()
        df_limited['ema_short'] = df_limited[self.base_column].ewm(span=self.ema_short_span).mean()
        df_limited['ema_long'] = df_limited[self.base_column].ewm(span=self.ema_long_span).mean()

        row = df_limited.iloc[-2]
        if row['ema_short']>row['ema_long']:
            ontop_prev = "s"
        else:
            ontop_prev = "l"
        ema_short_prev = row["ema_short"]
        ema_long_prev = row["ema_long"]
        # compare to now
        row = df_limited.iloc[-1]
        price = row[self.base_column]
        trade_type = None
        amount = 0

        if row['ema_short']>row['ema_long']:
            ontop_now = 's'
        else:
            ontop_now = 'l'
        self.ema_short_now = row["ema_short"]
        self.ema_long_now = row["ema_long"]
        self.price = price
        if ontop_now == 's' and ontop_prev == 'l':
            trade_type = "buy"
            decision = True
            amount = buy_bitcoin_with_USD(usd,trade_percentage,price)

        if ontop_now == 'l' and ontop_prev == 's':
            decision = True
            trade_type = "sell"
            amount = sell_bitcoin_for_USD(bitcoin,trade_percentage,price)
        properties = {
            "ontop_now":ontop_now,
            "ontop_prev":ontop_prev,
            "price":price,
            "ema_short_prev":ema_short_prev,
            "ema_long_prev":ema_long_prev

        }
        return decision, amount, trade_type, properties

def sell_bitcoin_for_USD(bitcoin,trade_percentage,price):
    how_many = ((bitcoin*trade_percentage)*price)
    return how_many

def buy_bitcoin_with_USD(usd,trade_percentage,price):
    how_many = usd*trade_percentage
    return how_many