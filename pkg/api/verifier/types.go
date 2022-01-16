package verifier

import "context"

type ITransactionVerifier interface {
	VerifyTransaction(ctx context.Context, req *VerifyTransactionRequest) (*VerifyTransactionResponse, error)
	Close()
}
