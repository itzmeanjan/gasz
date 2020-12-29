## Gas Price Fetcher

This simple Python script periodically sends HTTP GET request to ethgasstation.info & fetches latest gas price data, which is eventually published to a Redis PubSub channel. Now from there it's upto channel subscribers, how they want to process messages being received by them.
