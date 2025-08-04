# Project layout
```
.
├── cmd/server          # main.go
├── internal
│   ├── wallet          # domain: entity, repo interface, in-memory impl
│   ├── game            # publish GameOutcome events
│   ├── leaderboard     # publish Leaderboard events
│   └── transport
│       ├── http        # chi handlers & router
│       ├── ws          # hub + client adapter
│       └── grpc        # gRPC server adapter
├── pkg/api/wallet/v1       # generated pb.go files (buf)
├── proto/...           # .proto source + buf configs
├── Dockerfile
└── docker-compose.yml
```
# 📦 Selesa Slots — Wallet & Realtime Backend (Test Task)

Backend service for a gambling platform that

* keeps user wallets in an **in-memory** store (deposit / withdraw / balance),
* publishes **real-time** events via **WebSocket** (wallet, game outcomes, leaderboard),
* exposes a **gRPC** endpoint for `GetBalance`,
* ships in a single Docker image.

> **Stack**: 
> Go 1.23 · chi · gorilla/websocket · grpc-go · zap · buf · Docker multi-stage.

---

## ▶︎ Quick start (Docker)

```bash
# build & run in the background
make up

# REST demo — deposit 500 cents to user 1
curl -X POST http://localhost:8080/api/wallet/deposit \
     -H "Content-Type: application/json" \
     -d '{"user_id":1,"amount":500}'

# check balance via gRPC
grpcurl -plaintext -d '{"user_id":1}' \
        localhost:9091 \
        wallet.v1.WalletService/GetBalance
```
### Ports
```
Host	    Container	   Purpose
8080	    8080	   HTTP + WebSocket
9091	    9091	   gRPC
```

### REST API
```
|Method	|   Path	                     |   Body → Params	                    |   2xx response             | 
|----------------------------------------------------------------------------------------------------------------|
|POST	|   /api/wallet/deposit	             |   {user_id:int64, amount:int64 > 0}  |	 {user_id, balance}      |
|POST	|   /api/wallet/withdraw	     |   same	                            | same — 409 if insufficient |
|GET	|   /api/wallet/balance/{user_id}    |	 none	                            | {user_id, balance}         |
```
*Money is stored as integer “cents”, no floats.*



### WebSocket
```
GET /ws?topics=wallet,game,leaderboard
```
* text frames, one JSON per event
* ping/pong every 50 s keeps the connection alive
* back-pressure: if client’s send-buffer = 128 messages → connection closes

Example:
```json
{
  "type":"deposit",
  "user_id":1,
  "amount":500,
  "balance":500,
  "timestamp":"2025-08-04T03:59:58.422282263Z"
}
```
```json
{
  "game_id":"abc123",
  "user_id":7,
  "outcome":"win",
  "amount":2500,
  "timestamp":"2025-08-04T02:10:40.757131Z"
}
```
```json
{
  "user_id":7,
  "position":1,
  "score":9999,
  "updated_at":"2025-08-04T02:11:00.783514Z"
}
```
Publish test events (in another terminal):
```bash
# Game outcome
curl -X POST http://localhost:8081/api/game/outcome \
     -H "Content-Type: application/json" \
     -d '{"game_id":"abc", "user_id":7, "outcome":"win", "amount":2500}'

# Leaderboard update
curl -X POST http://localhost:8080/api/leaderboard/update \
     -H "Content-Type: application/json" \
     -d '{"user_id":7, "position":1, "score":9999}'

```

## gRPC
```protobuf
service WalletService {
  rpc GetBalance (GetBalanceRequest) returns (GetBalanceResponse);
}
```
With reflection enabled you can grpcurl without proto files:
```bash
grpcurl -plaintext -d '{"user_id":1}' \
        localhost:9091 \
        wallet.v1.WalletService/GetBalance
```

# Local development
```bash
go run ./cmd/server               # (or make run) default ports 8080/9091
go test ./... -race               # unit & integration
buf generate                      # re-generate protobuf stubs
```

### Notes & trade-offs
* In-memory store meets test-task spec; swap­pable via interface.
* Money = int64 (cents) to avoid FP errors.
* Concurrency-safe with sync.RWMutex.
* Simple /healthz for Docker HC; Prometheus & JWT skipped to stay within scope.