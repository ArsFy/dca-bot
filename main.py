import json
import time
import logging
from decimal import Decimal, ROUND_DOWN
from pybit.unified_trading import HTTP

# Logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

# Config
with open("config.json", "r") as f:
    config = json.load(f)

session = HTTP(
    demo=config.get('demo', False),
    api_key=config.get('api_key'),
    api_secret=config.get('api_secret'),
)

class DCA:
    def __init__(self, symbol, amount):
        self.symbol = symbol
        self.amount = Decimal(amount)
        self.price_gap = Decimal(2.4 * 0.01)
        self.price_scale = Decimal(1.7)
        self.amount_scale = Decimal(1.5)
        self.order_number = 6

        self.tickSize = Decimal(0.1)
        self.qtyStep = Decimal(0.1)

    # 獲取價格
    def get_price(self):
        data = session.get_tickers(
            category="linear",
            symbol=self.symbol
        )
        return Decimal(data['result']['list'][0]['lastPrice'])
    
    # 獲取合約
    def get_instruments_info(self):
        data = session.get_instruments_info(
            category="linear",
            symbol=self.symbol
        )
        data = data['result']['list'][0]
        self.tickSize = Decimal(data['priceFilter']['tickSize'])
        self.qtyStep = Decimal(data['lotSizeFilter']['qtyStep'])

    def place_order(self, price, qty):
        session.place_order(
            category='linear',
            symbol=self.symbol,
            side="Buy",
            orderType="Market",
            orderFilter="StopOrder",
            triggerPrice=str(price),
            triggerDirection="2",
            triggerBy="LastPrice",
            qty=str(qty),
            positionIdx="1",
            reduceOnly=False
        )

    def run(self):
        self.get_instruments_info()

        now_price = self.get_price()
        # now_price = Decimal(0.46240).quantize(self.tickSize, rounding=ROUND_DOWN)
        price_list = [
            now_price,
            (now_price - now_price * self.price_gap).quantize(self.tickSize, rounding=ROUND_DOWN)
        ]
        for index in range(1, self.order_number-1):
            last_price = price_list[-1]
            price_list.append(
                (last_price - last_price * self.price_gap * (self.price_scale ** index)).quantize(self.tickSize, rounding=ROUND_DOWN)
            )

        total_ratio = sum(self.amount_scale ** i for i in range(self.order_number))
        base_amount = self.amount / total_ratio
        amount_list = [((base_amount * (self.amount_scale ** i)) / price_list[i]).quantize(self.qtyStep, rounding=ROUND_DOWN) for i in range(self.order_number)]

        for i in range(1, self.order_number):
            self.place_order(price_list[i], amount_list[i])

        print(price_list, "\n", amount_list)
        

if __name__ == "__main__":
    bot1 = DCA("DOGEUSDT", 100)
    bot1.run()
