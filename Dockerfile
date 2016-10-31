FROM golang

ADD tg_nd_twtr_bot.go  /go/src/github.com/nektodev/tg_nd_twtr_bot/

RUN go get github.com/coreos/pkg/flagutil && go get github.com/dghubble/oauth1 && go get github.com/dghubble/go-twitter/twitter && go get gopkg.in/telegram-bot-api.v4 && go install github.com/nektodev/tg_nd_twtr_bot
ENTRYPOINT /go/bin/tg_nd_twtr_bot