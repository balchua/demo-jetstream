package verifier

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockTransactionVerifier struct {
	mock.Mock
}

func (m *MockTransactionVerifier) VerifyTransaction(ctx context.Context, req *VerifyTransactionRequest) (*VerifyTransactionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*VerifyTransactionResponse), args.Error(1)
}

func (m *MockTransactionVerifier) Close() {
	zap.S().Debug("mock close")
}
