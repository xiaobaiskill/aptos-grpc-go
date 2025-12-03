aptos-grpc-go
---

Golang code generated from [aptos-core/protos/proto](https://github.com/aptos-labs/aptos-core/tree/main/protos/proto)

### how to generate golang code for proto
* 1、Add go_package to proto files
Add the following option to each proto file:
```
option go_package = "github.com/xiaobaiskill/aptos-grpc-go/aptos/indexer/v1;indexerv1";
```

* 2、Create Makefile in `aptos-core/protos/proto`
```

OUT_DIR=../golang

.PHONY: generate
generate:
	protoc -I . \
      --go_out=$(OUT_DIR) \
      --go-grpc_out=$(OUT_DIR) \
      --go_opt=paths=source_relative \
      --go-grpc_opt=paths=source_relative \
      aptos/indexer/v1/*.proto aptos/internal/fullnode/v1/*.proto aptos/transaction/v1/*.proto aptos/util/timestamp/*.proto
```

* 3、genrate golang code
```
cd protos/proto
mkdir -p ../golang
make generate
```
