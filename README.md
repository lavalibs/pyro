# Pyro
\[WIP\] A scalable state container for Lavalink.

## Goals
- Complete coverage of Lavalink features
- Standalone usability as a simple Lavalink wrapper
- Additional utility features
  - Song queues
  - Clustering
- Horizontal scalability
  - Reliability through redundancy

## Architecture
Pyro is designed around 2 complimentary services. One directly faces each Lavalink node and provides an interface to Redis (we'll call these "Nodes"). The other acts as a server in front of Redis to essentially "load balance" the cluster by forwarding requests to appropriate nodes (we'll call these "Servers"); this service is necessary for three reasons:

1. **Reduce PubSub load on Redis:** VSUs and other player-specific queries (like play, pause, etc.) no longer need to broadcast on Redis; they can instead be sent directly to the node responsible for that player
2. **Improve consistency:** as a direct result of #1, the server is aware of which node it *should* be sending to and can take action in the case that the request fails somehow
3. **Reduce logic complexity:** figuring out how to create a player between a network of nodes would require complex consensus logic; in this architecture, the server can simply create a player on whatever node it deems appropriate

```
external stuff -> server -> Redis -> node -> Lavalink
```

Nodes are responsible for notifying Redis of their death, upon which other nodes in the cluster will automatically take over the dying node's players. If the node is killed before it can broadcast, players that it owned are dead until another node sweeps them up or the player is recreated manually; this situation is unlikely and should only occur in improperly configured environments.
