FROM alpine:3.15.4

RUN apk add go make

WORKDIR /rated

COPY . .

RUN make

CMD ["bin/rated-cli", "--config", "config.yaml", "watch"]
