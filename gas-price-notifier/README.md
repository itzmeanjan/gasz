## Gas Price Notifier

### Introduction

This simple backend application serves two purposes

- Static content required by `gasz` webUI
- Manages Ethereum gas price subscription/ unsubscription, over websocket connection
    - For ☝️ purpose, it connects to Redis PubSub topic, where gas prices are published by [gas-price-fetcher](../gas-price-fetcher) module and this application considers whether received price is what client is interested in or not

### Usage

- Make sure you've Golang _(>=1.15)_
- Make sure you've installed Redis and set up password based authentication.
- After cloning this repository, get inside this directory and run 👇 command, which downloads all dependencies for you

```bash
go get
```

- Create a `.env` file in this directory, with 👇 content

```
Port=7000
RedisHost=127.0.0.1
RedisPort=6379
RedisPassword=password
RedisDB=0
RedisPubSubChannel=gas-price
```

- You can build it now

```bash
go build
```

- Running 👇 starts backend application

```bash
./gas-price-notifier
```

> Make sure you've started `gas-price-fetcher` script before starting ☝️

> Note: For production, consider running it with `systemd`
