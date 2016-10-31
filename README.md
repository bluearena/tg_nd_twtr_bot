# Twitter-To-Telegram bot
Author Tsykin V.A. aka NektoDev
## Description

Simple telegram bot written on Golang that provides user twitter stream to Telegram. 

## Quickstart

1. Build container:
```
docker build -t tg_nd_twtr_bot:latest .
```

2. _(optional)_ Push container:
docker push tg_nd_twtr_bot:latest

3. Run container:
```
sudo docker run --restart=always --name=tg_nd_twtr_bot -d \
    -e TWITTER_BOT_CONSUMER_KEY= \
    -e TWITTER_BOT_CONSUMER_SECRET= \
    -e TWITTER_BOT_ACCESS_TOKEN= \
    -e TWITTER_BOT_ACCESS_SECRET= \
    -e TWITTER_BOT_TG_TOKEN= \
    -e TWITTER_BOT_TG_CHAT-ID=  \
    nektodev/tg_nd_twtr_bot
```

