package verifier

import (
	"context"
	"fmt"
	"log"

	"github.com/balchua/demo-jetstream/pkg/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TransactionVerifierClientWrapper struct {
	client TransactionVerifierClient
	conn   *grpc.ClientConn
}

func NewTransactionVerifierClientWrapper(config config.ApiConfig) (*TransactionVerifierClientWrapper, error) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := NewTransactionVerifierClient(conn)

	return &TransactionVerifierClientWrapper{
		client: c,
		conn:   conn}, nil
}

func (wrapper *TransactionVerifierClientWrapper) Close() {
	zap.S().Info("closing the connection to the api")
	wrapper.conn.Close()
}

func (wrapper *TransactionVerifierClientWrapper) VerifyTransaction(ctx context.Context, request *VerifyTransactionRequest) (*VerifyTransactionResponse, error) {
	response, err := wrapper.client.VerifyTransaction(ctx, request)
	return response, err

}
