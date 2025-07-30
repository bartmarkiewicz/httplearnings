package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaders_Parse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestHeaders_ParseCaseInsensitive(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestHeaders_ParseInvalidChar(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("host@: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeaders_ParseEmptyKey(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte(": localhost:42069\r\n\r\n")
	_, _, err := headers.Parse(data)
	require.Error(t, err)
}

func TestHeaders_ParseMultipleValsSameHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Set-Person: lane-loves-go\r\n")
	data2 := []byte("Set-Person: prime-loves-zig\r\n\r\n")
	headers.Parse(data)
	n, done, err := headers.Parse(data2)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 29, n)
	assert.False(t, done)
	assert.Equal(t, "lane-loves-go, prime-loves-zig", headers["set-person"])
}

func TestHeaders_Parse_InvalidSpacing(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
