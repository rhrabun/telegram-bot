import json

import boto3

from ..config import (
    AWS_REGION,
    AWS_ACCESS_KEY_ID,
    AWS_SECRET_ACCESS_KEY
)


AWS_SESSION = boto3.Session(
    region_name = AWS_REGION,
    aws_access_key_id=AWS_ACCESS_KEY_ID,
    aws_secret_access_key=AWS_SECRET_ACCESS_KEY
)


class AWSHelper(object):
    def __init__(self):
        pass

    @classmethod
    def get_scripts_list(self):
        client = AWS_SESSION.client('lambda')

        response = client.list_functions()
        functions = [func['FunctionName'] for func in response['Functions']]
        
        return functions

    @classmethod
    def run_function(self, func_name: str):
        client = AWS_SESSION.client('lambda')

        response = client.invoke(
            FunctionName=func_name,
            InvocationType='RequestResponse',
            LogType='None'
        )
        
        response_payload = json.loads(response['Payload'].read().decode("utf-8"))

        status_code = message = ''
        try:
            status_code = response_payload['statusCode']
            message = response_payload['message']
        except KeyError:
            status_code = 500
            message = response_payload['errorMessage']

        return status_code, message
