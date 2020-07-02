# go-wol-server

[![Build Status](https://cloud.drone.io/api/badges/alinbalutoiu/go-wol-server/status.svg)](https://cloud.drone.io/alinbalutoiu/go-wol-server)

Wake on LAN server written in go.

## How to use go-wol-server

The easiest way to get started is to use Docker to run it.

### Run the Docker container with go-wol-server

Start the server by running:
```
docker run -d --net=host --restart=always --name=gowolserver \
    -e PORT=8080 \
    -e LOG_LEVEL=info \
    alinbalutoiu/go-wol-server
```

If the server starts successfully, you should see the following:
```
$ docker logs gowolserver
time="2020-07-02T20:42:48Z" level=info msg="App initialized"
time="2020-07-02T20:42:48Z" level=info msg="Gorilla router initialized"
time="2020-07-02T20:42:48Z" level=info msg="Listening on port: 8080..."
```

The restart flag will make the container to start automatically on server reboot.
The easiest option is to run with the network mode set to `host`.

Supports the following environment variables:

- `PORT` is the port on which to listen (default: `8080`)

- `LOG_LEVEL` is the verbosity of the logs (default: `info`)

The Docker image is already built for the following architectures:
`arm`, `arm64` and `amd64`.

### Interacting with the server

Once the server is up and running, you can check the connectiong
by doing a `GET` request. For example using `curl`:

```
$ curl http://localhost:8080
Hello World!
```

If you don't get any output, make sure that the server is reachable
from your device and that there are no firewalls in between.

To instruct the server to send a wake-on-lan packet, do a `GET`
request at the endpoint `http://localhost:8080/wakeonlan/<mac_address>`
where `mac_address` is a IEEE 802 MAC-48 address.

Examples:
```
$ curl http://localhost:8080/wakeonlan/11:22:33:44:55:66
{"message": "OK"}
$ curl http://localhost:8080/wakeonlan/11-22-33-44-55-66
{"message": "OK"}
$ curl http://localhost:8080/wakeonlan/11-22-33-44-55
{"message": "address 11-22-33-44-55: invalid MAC address"}
```

HTTP status code `200` is returned if the request is completed
successfully, and `400` otherwise.
