# gasz

Ethereum Gas Price Notifier

![banner](sc/banner.gif)

## Introduction

While interacting with Ethereum Network I had to frequently visit `ethgasstation.info` for checking current `safeLow` gas price and decide when I can send transaction, which was really eating a lot of my time.

So I started looking for some exisiting solutions which can give me Ethereum Gas Price feed, but couldn't find anything that satisfies my requirement. That's when I decided to write my own solution for getting notified when Ethereum Gas Price reaches a certain threshold.

After some exploration, I decided to build this application while leveraging PubSub model, because of 👇

- One module of `gasz` keeps fetching latest Ethereum Gas Price from `ethgasstation.info` [ **gas-price-fetcher** ]
- It also publishes latest gas prices to a Redis PubSub topic
- Another module accepts client request for gas price reaching certain threshold i.e. _<= 50Gwei_, over Websocket [ **gas-price-notifier** ]
- Client will subscribe to Redis PubSub topic & will receive notification over Websocket, when gas price reaches that value
- Single client can subscribe to multiple gas price feeds over same Websocket connection 🦾

This worked pretty well. After basic building block, I built a very minimalistic webUI for subscription/ unsubscription management & displaying HTML5 Notification.

And here's `gasz`.

> Note: `gasz` will be live very soon.

**More coming soon ...**
