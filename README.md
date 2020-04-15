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

- <strong>server</strong>
```
go run cli/main.go -mode=server -addr=0.0.0.0
```

- <strong>client1</strong>
```
go run cli/main.go -mode=client -addr=0.0.0.0:9001
// random port will generate for client1
// e.g listening on [::]:44246
```

- <strong>client2</strong> connect to <strong>client1</strong>
```
go run cli/main.go -mode=client -addr=0.0.0.0:44246
// random port will generate for client2
// e.g listening on [::]:41007
```
