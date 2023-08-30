Service for VerifyingPaymaster
==============================

## RPC

```
curl -X POST http://localhost:8888/rpc/1234567890 -H "Content-Type:application/json" --data '{
    "jsonrpc":"2.0",
                "method":"pm_sponsorUserOperation",
                "params":[],
    "id":1
}'

curl -X POST http://localhost:8888/rpc/1234567890 -H "Content-Type:application/json" --data '{
    "jsonrpc":"2.0",
                "method":"pm_gasRemain",
                "params":["0x816117a3E3A909947e9835d3904A2991696F1FD2"],
    "id":1
}'
```

## Docker

```
docker build -t paymaster:latest .
```
