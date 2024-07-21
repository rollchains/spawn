# Nginx Explorer

To run the explorer locally, this nginx configuration is required to get https to work locally.

## PingPub Docker Build

```dockerfile
# docker build . -t pingpub:latest

FROM node:20-alpine

RUN apk add --no-cache yarn

WORKDIR /app

COPY . .

EXPOSE 8080

CMD [ "yarn", "--ignore-engines", "serve", "--host", "0.0.0.0" ]
```

## Running

Update your `/etc/hosts` file to include the following:

```
127.0.0.1 api.localhost
127.0.0.1 rpc.localhost
127.0.0.1 pingpub.localhost
```

Then `docker compose up` to start the reverse proxy, explorer, and the RPC/REST API Services. Then start the testnet (make sh-testnet)

