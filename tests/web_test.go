package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWebIsAlive(t *testing.T) {
	// TODO: this will be broken if the web container port changes
	resp, err := http.Get("http://localhost:6543")
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// TODO: more tests when needed!
