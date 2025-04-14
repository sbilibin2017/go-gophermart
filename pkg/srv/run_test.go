package srv

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun_SuccessfulStartAndShutdown(t *testing.T) {

	mockServer := &http.Server{
		Addr: "localhost:8080",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := Run(ctx, mockServer)
		require.NoError(t, err, "Server should shut down gracefully")
	}()

	time.Sleep(100 * time.Millisecond)

	cancel()

}
