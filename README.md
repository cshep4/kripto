# kripto ₿ 💰 💸 🤑 🏴󠁧󠁢󠁷󠁬󠁳󠁿

[![CircleCI](https://circleci.com/gh/cshep4/kripto.svg?circle-token=86c9f9b058b912c8b87271abf4f054c5ce9451a5)](https://circleci.com/gh/cshep4/kripto)

Kripto ₿ trading platform, periodically checks the BTC-GBP exchange rate using Coinbase APIs and makes intelligent decisions whether to buy/sell Bitcoin. Executes trades using Coinbase Pro.

## Functions λ

| Function                                                | Service                                     | Runtime       | Events             | Description                                                                       |
| ------------------------------------------------------- | ------------------------------------------- | ------------- | ------------------ | --------------------------------------------------------------------------------- |
| [rate-retriever](./services/rate-retriever)             | [rate-retriever](./services/rate-retriever) | Node.js       | Schedule           | Retrieves the BTC-GBP exchange rate from Coinbase and publishes result to SQS.    |
| [trade](./services/trader/cmd/trade)                    | [trader](./services/trader)                 | Go            | Invocation         | Calls Coinbase Pro to make a BTC-GBP trade and publishes result to SQS.           |
| [rate-writer](./services/data-storer/cmd/rate-writer)   | [data-storer](./services/data-storer)       | Go            | SQS                | Stores a trade in the database.                                                   |
| [trade-writer](./services/data-storer/cmd/trade-writer) | [data-storer](./services/data-storer)       | Go            | SQS                | Stores a rate in the database.                                                    |

### Rate Retriever 💲

- **Language** - JavaScript
- **Runtime** - nodejs12.x
- **Event** - Scheduled - every minute
- **Services** - AWS Lambda, Serverless, SQS (Producer), Coinbase API

##### Request
    {}

##### Response 
    {}
    
### Trade 🤝

- **Language** - Go
- **Runtime** - go1.x
- **Event** - Invocation
- **Services** - AWS Lambda, Serverless, SQS (Producer), Coinbase Pro API

##### Request
    {
        "tradeType": "buy"
    }

##### Response 
    {}

### Rate Writer 💰

- **Language** - Go
- **Runtime** - go1.x
- **Event** - SQS - `RateUpdate` queue
- **Services** - AWS Lambda, Serverless, SQS (Consumer), MongoDB
- **Idempotency** - `idempotencyKey` sent in message payload

##### Request
    {
        "idempotencyKey": "aa368788-bb4f-40c0-b80f-afcfdaf18574",
        "rate": "8012.92",
        "dateTime": "2020-05-19T19:39:00"
    }

##### Response 
    {}

### Trade Writer 💸

- **Language** - Go
- **Runtime** - go1.x
- **Event** - SQS - `Trade` queue
- **Services** - AWS Lambda, Serverless, SQS (Consumer), MongoDB
- **Idempotency** - trade ID (`id`) in message payload used as idempotency key

##### Request
    {
        "id": "aa368788-bb4f-40c0-b80f-afcfdaf18574",
        "side": "buy",
        "productId": "BTC-GBP",
        "funds": "9.95024875",
        "settled": true,
        "createdAt": "2020-05-19T19:39:00",
        "fillFees": "0.049751102976",
        "filledSize": "0.00125952",
        "executedValue": "9.9502205952"
    }

##### Response 
    {}
    
## Queues ✉️

### Trade

- **Description** - Signifies a trade event has taken place
- **Producers** - trader
- **Consumers** - trade-writer

##### Payload
    {
        "id": "aa368788-bb4f-40c0-b80f-afcfdaf18574",
        "side": "buy",
        "productId": "BTC-GBP",
        "funds": "9.95024875",
        "settled": true,
        "createdAt": "2020-05-19T19:39:00",
        "fillFees": "0.049751102976",
        "filledSize": "0.00125952",
        "executedValue": "9.9502205952"
    }
    
### RateUpdate

- **Description** - Signifies a rate update event
- **Producers** - rate-retriever
- **Consumers** - rate-writer

##### Payload
    {
        "idempotencyKey": "aa368788-bb4f-40c0-b80f-afcfdaf18574",
        "rate": "8012.92",
        "dateTime": "2020-05-19T19:39:00"
    }