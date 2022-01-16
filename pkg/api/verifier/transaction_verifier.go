package verifier

import (
	context "context"

	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type TransactionVerifierImpl struct {
	UnimplementedTransactionVerifierServer
}

func (tv *TransactionVerifierImpl) VerifyTransaction(ctx context.Context, req *VerifyTransactionRequest) (*VerifyTransactionResponse, error) {
	_, span := otel.Tracer("verifier").Start(ctx, "RemoteVerifyTransaction")
	defer span.End()
	zap.S().Infof("transaction ID %s", req.Tx.TransactionID)

	amountInDecimal, err := decimal.NewFromString(req.Tx.Amount)
	if err != nil {
		return nil, err
	}
	zap.S().Debugf("Transaction request amount %v", amountInDecimal)
	if amountInDecimal.LessThanOrEqual(decimal.Zero) {
		return &VerifyTransactionResponse{
			Tx:      req.Tx,
			Message: "Invalid amount sent",
			Code:    StatusCode_NOT_OK,
		}, nil
	}
	return &VerifyTransactionResponse{
		Tx:      req.Tx,
		Message: "Valid Transaction",
		Code:    StatusCode_OK,
	}, nil

}
