Service for VerifyingPaymaster
==============================

## RPC

```
curl -X POST http://localhost:8888 -H "Content-Type:application/json" --data '{
    "jsonrpc":"2.0",
                "method":"eth_signVerifyingPaymaster",
                "params":[],
    "id":1
}'
```
