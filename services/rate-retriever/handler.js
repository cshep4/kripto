const pino = require('pino');
const logger = pino({
    name: 'rate-retriever',
    messageKey: 'message',
    changeLevelName: 'severity',
    useLevelLabels: true
});

logger.info("initialisation");

const AWS = require('aws-sdk');
// const sqs = new AWS.SQS({
//     region: process.env.REGION
// });
const sns = new AWS.SNS({
    region: process.env.REGION
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
    // const queueUrl = process.env.QUEUE_URL;
    const topic = process.env.TOPIC;

    coinbaseClient.getBuyPrice({'currencyPair': 'BTC-GBP'}, function (err, price) {
        if (err) {
            logger.error({"msg": "failed to send message", "currencyPair": "BTC-GBP", "err": err});
            callback(err);
            return;
        }

        // const params = {
        //     MessageBody: JSON.stringify({
        //         rate: parseFloat(price.data.amount),
        //         dateTime: new Date(),
        //         idempotencyKey: uuidv4(),
        //     }),
        //     QueueUrl: queueUrl
        // };

        let params = {
            Message: JSON.stringify({
                rate: parseFloat(price.data.amount),
                dateTime: new Date(),
                idempotencyKey: uuidv4(),
            }),
            TopicArn: topic
        };

        sns.publish(params, function(err, data) {
            if (err) {
                logger.error({"msg": "failed to publish event", "body": params.MessageBody, "err": err});
                callback(err);
            }
        });

        // sqs.sendMessage(params, function (err, data) {
        //     if (err) {
        //         logger.error({"msg": "failed to send message", "body": params.MessageBody, "err": err});
        //         callback(err);
        //     }
        // });
    });
});