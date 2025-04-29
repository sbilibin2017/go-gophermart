package services

import (
	"testing"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapToStruct_Success(t *testing.T) {
	// Prepare the data for the test
	data := map[string]any{
		"Number":     "9278923470",
		"Accrual":    500, // This will map to *int64
		"Status":     "PROCESSED",
		"UploadedAt": "2020-12-10T15:15:45+03:00", // Time in RFC3339 format
	}

	// Initialize an empty Order struct
	var order domain.Order

	// Call the mapToStruct function to map the data to the struct
	err := mapToStruct(&order, data)

	// Assert that there is no error
	require.NoError(t, err)

	// Assert that the fields in the struct are correctly populated
	assert.Equal(t, "9278923470", order.Number)
	assert.Equal(t, int64(500), *order.Accrual) // Assert for *int64 field
	assert.Equal(t, "PROCESSED", order.Status)

	// Validate the time is correctly parsed and matches the expected time
	expectedTime := time.Date(2020, 12, 10, 15, 15, 45, 0, time.FixedZone("UTC+3", 3*60*60))
	assert.Equal(t, expectedTime, order.UploadedAt)
}
