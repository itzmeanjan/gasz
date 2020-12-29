#!/usr/bin/python3

from typing import Any, Dict
from requests import get
from requests.exceptions import Timeout
from requests.exceptions import HTTPError
import logging

logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s')

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



if __name__ == '__main__':
    print('[!] This module is not supposed to be used this way !')
