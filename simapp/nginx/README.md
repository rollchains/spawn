# Nginx Explorer

To run the explorer locally, this nginx configuration is required to get local https to work.

## Running

Update your `/etc/hosts` file to include the following:

```
127.0.0.1 api.localhost
127.0.0.1 rpc.localhost
127.0.0.1 pingpub.localhost
```

Start the testnet with: `make sh-testnet` or the full IBC network with `make testnet`

Then `docker compose up` to start the reverse proxy, explorer, and the RPC/REST API Services.

<!-- markdown-link-check-disable-next-line -->
Visit: https://pingpub.localhost to view the explorer.

> Attempting to view as a standard http:// instance will break the block explorer due to pesky CORS errors.
