# kripto ₿ 💰 💸 🤑

[![CircleCI](https://circleci.com/gh/cshep4/kripto.svg?circle-token=86c9f9b058b912c8b87271abf4f054c5ce9451a5)](https://circleci.com/gh/cshep4/kripto)

## REST API

The REST API to the Kripto ₿ platform is described below.

### Get Historic Rates & Trades

#### Request

`GET /data`

    curl -H "Content-Type: application/json" -X GET /data

#### Response

    HTTP/1.1 200 OK
    Date: Mon, 14 May 2020 22:10:00 GMT
    Status: 200 OK
    Access-Control-Allow-Origin: *
    Content-Type: application/json

    {
        "rates": [{
            "rate": "7879.4087",
            "dateTime": "2020-05-14T22:18:00"
        }],
        "trades": [{
            "value": {
                "gbp": 7879.41,
                "btc": 1,
            },
            "type": "buy",
            "rate": "7879.4087",
            "dateTime": "2020-05-14T22:18:00"
        }]
    }

### Store rate & trade

#### Request

`POST /data`

    curl -d '{"rate": { "rate": "7879.4087", "dateTime": "2020-05-14T22:18:00" }, "buy": { traded: true, "gbp": 7879.41, "btc": 1, "rate": "7879.4087", "dateTime": "2020-05-14T22:18:00" }, "sell": { traded: true, "gbp": 7879.41, "btc": 1, "rate": "7879.4087", "dateTime": "2020-05-14T22:18:00" } }' -H "Content-Type: application/json" -X POST /data

###### Request Body

    {
        "rate": {
            "rate": "7879.4087",
            "dateTime": "2020-05-14T22:18:00"
        },
        "buy": {
            traded: true,
            "gbp": 7879.41,
            "btc": 1,
            "rate": "7879.4087",
            "dateTime": "2020-05-14T22:18:00"
        },
        "sell": {
            traded: true,
            "gbp": 7879.41,
            "btc": 1,
            "rate": "7879.4087",
            "dateTime": "2020-05-14T22:18:00"
        }
    }

#### Response

    HTTP/1.1 201 Created
    Date: Mon, 14 May 2020 22:10:00 GMT
    Status: 201 Created
    Access-Control-Allow-Origin: *
    Content-Type: application/json

    {}
