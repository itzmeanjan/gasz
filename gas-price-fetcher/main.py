#!/usr/bin/python3

from dotenv import load_dotenv
from os import getenv
from typing import Any, Dict
from requests import get
from requests.exceptions import Timeout
from requests.exceptions import HTTPError
from redis import Redis

import logging
logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s', level=logging.DEBUG)

load_dotenv()

GasPriceProducer = getenv('GasPriceProducer')
RedisHost = getenv('RedisHost')
RedisPort = getenv('RedisPort')
RedisPassword = getenv('RedisPassword')


def connectToRedis() -> Redis:
    '''
        Connecting to Redis instance and returning
        connection object back to caller
    '''
    conn = None
    try:
        conn = Redis(host=RedisHost, port=RedisPort, password=RedisPassword)
        logging.info('Connected to Redis')
    except Exception as e:
        logging.error(f'{e}')
    finally:
        return conn


def fetchGasPrice(url: str) -> Dict[str, Any]:
    '''
        Given URL of end point producing gas price feed,
        it'll query that endpoint and return JSON response back
        to function caller
    '''
    data = {}
    try:
        resp = get(url, timeout=1)

        if not (resp.status_code == 200):
            raise Exception('Response with non-200 status code')
        
        data = resp.json()
        resp.close()
    except HTTPError as e:
        logging.error(f'{e}')
    except Timeout as e:
        logging.error(f'{e}')
    except Exception as e:
        logging.error(f'{e}')
    finally:
        return data

def main():
    conn = connectToRedis()
    print(conn)
    print(fetchGasPrice(GasPriceProducer))

if __name__ == '__main__':
    main()
