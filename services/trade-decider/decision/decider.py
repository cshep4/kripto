import logging
import pandas as pd


class Decider:
    def __init__(self, logger: logging.Logger):
        self.logger = logger
        self.rolling_window = 100
        self.base_column = "p"  # price/rate
        self.ema_short_span = 9
        self.ema_long_span = 21
        self.trade_percentage = 0.1

    def decide(
            self,
            data: dict,
            bitcoin: float,
            gbp: float
    ) -> (bool, float, str, dict):
        decision = False

        df_limited = pd.DataFrame.from_dict(data)
        df_limited['rolling'] = df_limited[self.base_column].rolling(self.rolling_window).mean()
        df_limited['ema_short'] = df_limited['rolling'].ewm(span=self.ema_short_span).mean()
        df_limited['ema_long'] = df_limited['rolling'].ewm(span=self.ema_long_span).mean()

        row = df_limited.iloc[-2]

        if row['ema_short'] > row['ema_long']:
            ontop_prev = "s"
        else:
            ontop_prev = "l"

        # compare to now
        row = df_limited.iloc[-1]
        price = row["p"]
        trade_type = None
        amount = 0

        if row['ema_short'] > row['ema_long']:
            ontop_now = 's'
        else:
            ontop_now = 'l'

        if ontop_now == 's' and ontop_prev == 'l':
            trade_type = "buy"
            decision = True
            amount = buy_bitcoin(gbp, self.trade_percentage, price)

        if ontop_now == 'l' and ontop_prev == 's':
            decision = True
            trade_type = "sell"
            amount = buy_gbp(bitcoin, self.trade_percentage, price)

        return decision, amount, trade_type


def buy_gbp(bitcoin, trade_percentage, price):
    how_many = ((bitcoin * trade_percentage) * price)
    return how_many


def buy_bitcoin(gbp, trade_percentage, price):
    how_many = ((gbp * trade_percentage) / price)
    return how_many
