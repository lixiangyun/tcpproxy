# tcpproxy
- Tcp protocol proxy tool

## build
```
go build .
```

## usage
```
tcpproxy version: tcpproxy/0.1.0
Usage: tcpproxy [-h] [-m bridge/link/proxy] [-pools num] [-local ip:port] [-remote ip:port]

Options:
  -h    this help
  -local string
        connect to local address.
  -m string
        using bridge/link/proxy mode. (default "proxy")
  -pools uint
        using connect num on link/bridge mode. (default 10)
  -remote string
        connect to remote address.
```

## binary deployment
### proxy mode
```
tcpproxy.exe -m proxy -local 127.0.0.1:1000 -remote 10.10.0.1:2000
```

### bridge mode
```
tcpproxy.exe -m bridge -local :1000 -remote :2000
```

### link mode
```
tcpproxy.exe -m link -local 127.0.0.1:1000 -remote 10.10.0.1:2000
```

## Docker deployment
### proxy mode
```
docker run --net=host -d --restart=always linimbus/tcpproxy -m proxy -local 127.0.0.1:1000 -remote 10.10.0.1:2000
```

### bridge mode
```
docker run --net=host -d --restart=always linimbus/tcpproxy -m bridge -local :1000 -remote :2000
```

### link mode
```
docker run --net=host -d --restart=always linimbus/tcpproxy -m link -local 127.0.0.1:1000 -remote 10.10.0.1:2000
```
