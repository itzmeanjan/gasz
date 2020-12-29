## Gas Price Fetcher

### Introduction

This simple Python script periodically sends HTTP GET request to ethgasstation.info & fetches latest gas price data, which is eventually published to a Redis PubSub channel. Now from there it's upto channel subscribers, how they want to process messages being received by them.

### Usage

- Make sure you've Python _(>=3.7)_, Pip installed
- Make sure you've installed Redis and set up password based authentication.
- After cloning this repository, get inside this directory and run 👇 command, which downloads all dependencies for you

```bash
pip install -r requirements.txt
```

- Create a `.env` file in this directory, with 👇 content

```
GasPriceProducer=https://ethgasstation.info/api/ethgasAPI.json
RedisHost=127.0.0.1
RedisPort=6379
RedisPassword=password
RedisDB=0
RedisPubSubChannel=gas-price
SleepPeriod=2
RequestTimeout=2
```

- Now you can run 👇 from this directory

```bash
python3 main.py
```

> Note: For production, consider running it with `systemd`
