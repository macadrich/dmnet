[![Build Status](https://travis-ci.org/macadrich/dmnet.svg?branch=master)](https://travis-ci.org/macadrich/dmnet)
# dmnet
P2P protocol implemented in Golang

## Quick Start

0. Clone the repo

```
git clone https://github.com/macadrich/dmnet.git
cd dmnet
```

1. Install Go dependencies

```
dep ensure
```
2. Execute command

- server
```
go run cli/main.go -mode=server -addr=0.0.0.0
```

- client
```
go run cli/main.go -mode=client -addr=0.0.0.0
```
