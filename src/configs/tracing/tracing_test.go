package tracing_test

import (
	"context"
	"testing"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	tracing "github.com/KaueTTS/streaming_api/src/configs/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestInit(t *testing.T) {
	t.Run("should be successfully initialized", func(t *testing.T) {
		oldOTLPExporterEndpoint := env.OTLPExporterEndpoint
		oldAppName := env.AppName

		env.OTLPExporterEndpoint = "localhost:4318"
		env.AppName = "streaming_api_test"

		t.Cleanup(func() {
			env.OTLPExporterEndpoint = oldOTLPExporterEndpoint
			env.AppName = oldAppName
		})

		ctx := context.Background()

		tracerProvider, err := tracing.Init(ctx)

		require.NoError(t, err)
		require.NotNil(t, tracerProvider)

		t.Cleanup(func() {
			require.NoError(t, tracerProvider.Shutdown(ctx))
		})

		assert.Same(t, tracerProvider, otel.GetTracerProvider())
	})
}
