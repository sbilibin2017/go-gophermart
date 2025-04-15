package runners

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Init(zapcore.InfoLevel)
}

func TestRunServer_SuccessfulStartAndShutdown(t *testing.T) {

	mockServer := &http.Server{
		Addr: "localhost:8080",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := RunServer(ctx, mockServer)
		require.NoError(t, err, "Server should shut down gracefully")
	}()

	time.Sleep(100 * time.Millisecond)

	cancel()

}
