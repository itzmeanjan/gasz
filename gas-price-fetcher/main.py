#!/usr/bin/python3

from dotenv import load_dotenv
from os import getenv
from time import sleep
from typing import Any, Dict
from requests import get
from requests.exceptions import Timeout
from requests.exceptions import HTTPError
from redis import Redis
from json import dumps

import logging
# Setting up logging
logging.basicConfig(
    format='%(asctime)s - %(levelname)s - %(message)s', level=logging.DEBUG)

# Loading .env file content
#
# `.env` file is supposed to be present in this directory
load_dotenv()

# -- Reading configuration parameters
GasPriceProducer = getenv('GasPriceProducer')

RedisHost = getenv('RedisHost')
RedisPort = int(getenv('RedisPort'))
RedisPassword = getenv('RedisPassword')
RedisDB = getenv('RedisDB')
RedisPubSubChannel = getenv('RedisPubSubChannel')

SleepPeriod = int(getenv('SleepPeriod'))
RequestTimeout = int(getenv('RequestTimeout'))
# -- Done reading configuration paramaters


def connectToRedis() -> Redis:
    '''
        Connecting to Redis instance and returning
        connection object back to caller
    '''
    conn = None
    try:
        conn = Redis(host=RedisHost, port=RedisPort,
                     password=RedisPassword, db=RedisDB)
        conn.ping()

        logging.info('Connected to Redis')
    except Exception as e:
        conn = None
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
        resp = get(url, timeout=RequestTimeout)

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
    if not conn:
        exit(1)

    data = dict()
    while True:
        tmp = fetchGasPrice(GasPriceProducer)
        # If empty response received from ethgasstation
        if not tmp:
            sleep(SleepPeriod)
            continue

        # If previous data and current data same, no need to publish
        if data == tmp:
            sleep(SleepPeriod)
            continue

        # Only if data received gets updated ( i.e. different than previous iteration ),
        # then we'll publish it to channel
        #
        # Also note, as of now, `fast`, `fastest`, `safeLow`, `average` gas prices are
        # going to be published
        #
        # We can change it in future, if we're interested in providing more features
        # from listener side
        data = tmp
        logging.info(
            f'Published to {conn.publish(RedisPubSubChannel, dumps(dict([(i, data[i]/10) for i in ["fast", "fastest", "safeLow", "average"]])))} channel(s)')
        sleep(SleepPeriod)


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        logging.info('Killing program')
