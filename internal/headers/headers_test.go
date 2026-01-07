package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewHeaders() Headers {
	return make(map[string]string)
}

func TestParseHeader(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("HOSt: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra spaces
	headers = NewHeaders()
	data = []byte("   Host: localhost:42069   \r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["Host"] = "localhost:0"
	data = []byte("Host: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "localhost:42069", headers["host"])
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid header with uppercase
	headers = NewHeaders()
	data = []byte("HoST: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid header with unsupported char
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Joined 3 headers with same key
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go\r\n")
	n, done, err = headers.Parse(data)
	data = []byte("Set-Person: prime-loves-zig\r\n")
	n, done, err = headers.Parse(data)
	data = []byte("Set-Person: tj-loves-ocaml\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
	assert.False(t, done)
}
