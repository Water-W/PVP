# PVP

a P2P network Visualized Profiling tool

## demo

TODO

## install

```
go get -u github.com/Water-W/PVP/cmd/pvp
```

## usage 

master
```
./pvp -l 8000 -p 8001 # master listens at 8000 and HTTP port is 8001
```

worker
```
./pvp -m 1.2.3.4:8000 # worker connects master with given ip:port
```

