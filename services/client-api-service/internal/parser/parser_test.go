package parser_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/zsoltggs/golang-example/services/client-api-service/internal/parser"
	"github.com/zsoltggs/golang-example/services/client-api-service/internal/parser/mockpds"
	"os"
	"testing"
)

func TestParser_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPds := mockpds.NewMockServiceClient(ctrl)

	svc := parser.New(mockPds)
	ctx := context.Background()
	f, err := os.Open("test_data.json")
	require.NoError(t, err)

	mockPds.EXPECT().
		UpsertPort(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)

	err = svc.Parse(ctx, f)
	require.NoError(t, err)
}

// TODO Test unhappy path
// TODO Test exact input
