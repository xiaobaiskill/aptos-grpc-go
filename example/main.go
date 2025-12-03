package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"

	indexerV1 "github.com/xiaobaiskill/aptos-grpc-go/aptos/indexer/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	target = "grpc.mainnet.aptoslabs.com:443" // https://geomi.dev/
	token  = os.Getenv("APTOS_GRPC_TOKEN")    // export APTOS_GRPC_TOKEN=""
)

type authTokenCredential struct {
	token string
}

func NewAuthTokenCredential(token string) *authTokenCredential {
	return &authTokenCredential{token: token}
}

func (a authTokenCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + a.token,
	}, nil
}

func (a authTokenCredential) RequireTransportSecurity() bool {
	return true
}

func main() {
	client, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(NewAuthTokenCredential(token)),
	)
	if err != nil {
		log.Fatal("connect failed: ", err)
	}
	defer client.Close()

	dataClient := indexerV1.NewDataServiceClient(client)
	var (
		account = "0xc7efb4076dbe143cbcd98cfaaa929ecfc8f299203dfff63b95ccb6bfe19850fa"
		module  = "swap"
	)

	transactionRes, err := dataClient.GetTransactions(context.Background(), &indexerV1.GetTransactionsRequest{
		TransactionFilter: &indexerV1.BooleanTransactionFilter{
			Filter: &indexerV1.BooleanTransactionFilter_LogicalOr{
				LogicalOr: &indexerV1.LogicalOrFilters{
					Filters: []*indexerV1.BooleanTransactionFilter{
						{
							Filter: &indexerV1.BooleanTransactionFilter_ApiFilter{
								ApiFilter: &indexerV1.APIFilter{
									Filter: &indexerV1.APIFilter_EventFilter{
										EventFilter: &indexerV1.EventFilter{
											StructType: &indexerV1.MoveStructTagFilter{
												Address: &account,
												Module:  &module,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Fatal("get transactions failed: ", err)
	}
	defer transactionRes.CloseSend()
	var i = 0
	for {
		i++
		fmt.Println(i)
		if i > 5 {
			break
		}
		data, err := transactionRes.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("receive failed: ", err) // Failed to receive data, needs to subscribe again
			break
		}
		log.Println(data.GetProcessedRange().GetLastVersion())
	}
}
