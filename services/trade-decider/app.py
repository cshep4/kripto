import logging
# import os
import boto3
import json

# from lumigo_tracer import lumigo_tracer

# Set up logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)

logger.info('initialisation')

dynamo_client = boto3.client('dynamodb')
lambda_client = boto3.client('lambda')


# @lumigo_tracer(token=os.environ['LUMIGO_TRACER_TOKEN'], enhance_print=True)
def handler(event, context):
    logger.debug('Received event: {}'.format(event))

    resp = lambda_client.invoke(
        FunctionName="data-reader",
        InvocationType='RequestResponse',
        Payload="")

    data = json.loads(resp['Payload'].read())
    logger.info('{} data-reader success'.format(data))

    # for record in event['Records']:
    #     payload = json.loads(record['body'], parse_float=str)
    #     operation = record['messageAttributes']['Method']['stringValue']
    #     if operation in operations:
    #         try:
    #             operations[operation](dynamo_client, payload)
    #             logger.info('{} method successful'.format(operation))
    #             logger.debug('Event payload: {}'.format(payload))
    #         except Exception as e:
    #             logger.error(e)
    #     else:
    #         logger.error('Unsupported method \'{}\''.format(operation))
    return 0
