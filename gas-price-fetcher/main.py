#!/usr/bin/python3

from dotenv import load_dotenv
from os import getenv
import logging
from .fetch import fetchGasPrice

load_dotenv()

GasPriceProducer = getenv('GasPriceProducer')

logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s')

def main():
    logging.info(fetchGasPrice(GasPriceProducer))

if __name__ == '__main__':
    main()
