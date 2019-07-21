package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

// WriteJSONTests tests for WriteJSON
type WriteJSONTests struct {
	suite.Suite
}

// TestWrite
func (ts *WriteJSONTests) TestWrite() {
	// write simple boolean
	w := httptest.NewRecorder()
	MustWriteJSON(w, false, 200)
	ts.Equal(200, w.Code)
	ts.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))
	ts.Equal("false\n", w.Body.String())

	// write simple object
	w = httptest.NewRecorder()
	err := WriteJSON(w, map[string]interface{}{"message": "hello"}, 402)
	ts.NoError(err)
	ts.Equal(402, w.Code)
	ts.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))
	ts.Equal(`{"message":"hello"}`+"\n", w.Body.String())
}

// TestWriteJSON tests the WriteJSON
func TestWriteJSON(t *testing.T) {
	suite.Run(t, &WriteJSONTests{})
}
