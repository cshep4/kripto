const AWS = require('aws-sdk');
const sqs = new AWS.SQS({
    region: 'us-east-1'
});

const Client = require('coinbase').Client;
const coinbaseClient = new Client({
    'apiKey': process.env.COINBASE_API_KEY,
    'apiSecret': process.env.COINBASE_SECRET_KEY,
    strictSSL: false
});

const {v4: uuidv4} = require('uuid');

const lumigo = require('@lumigo/tracer')({token: process.env.LUMIGO_TRACER_TOKEN});

exports.handler = lumigo.trace((event, context, callback) => {
    const queueUrl = process.env.QUEUE_URL;

    coinbaseClient.getBuyPrice({'currencyPair': 'BTC-GBP'}, function (err, price) {
        if (err) {
            console.log('error:', "failed to get rate" + err);
            callback(err);
            return;
        }

        const params = {
            MessageBody: JSON.stringify({
                rate: parseFloat(price.data.amount),
                dateTime: new Date(),
                idempotencyKey: uuidv4(),
            }),
            QueueUrl: queueUrl
        };

        sqs.sendMessage(params, function (err, data) {
            if (err) {
                console.log('error:', "failed to send message" + err);
                callback(err);
                return;
            }

            console.log('data:', data.MessageId);
        });
    });
});