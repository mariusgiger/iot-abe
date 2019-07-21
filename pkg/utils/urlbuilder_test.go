package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// URLBuilderTests tests for URLBuilder
type URLBuilderTests struct {
	suite.Suite
}

// TestURLBuilder tests the URLBuilder
func TestURLBuilder(t *testing.T) {
	suite.Run(t, &URLBuilderTests{})
}

// TestParse tests for parse functions
func (ts *URLBuilderTests) TestParse() {
	// parse invalid URL
	b1, err := ParseURL("://")
	if ts.Error(err) {
		ts.Contains(err.Error(), "missing protocol scheme")
	}
	ts.Nil(b1)
	ts.Panics(func() {
		MustParseURL("://")
	})

	// parse valid URL
	b2, err := ParseURL("http://")
	ts.NoError(err)
	ts.NotNil(b2)
	ts.NotPanics(func() {
		MustParseURL("http://")
	})
}

// TestPath tests path modifications
func (ts *URLBuilderTests) TestPath() {
	// create URL builder
	b1 := func() *URLBuilder {
		return MustParseURL("http:///?foo=1")
	}

	// set path
	ts.Equal("http:///?foo=1", b1().String())
	ts.Equal("http:///info?foo=1", b1().SetPath("/info").String())
	ts.Equal("http:///info?foo=1", b1().SetPath("/info/").String())
	ts.Equal("http:///info/about?foo=1", b1().SetPath("/info/about").String())
	ts.Equal("http:///info/about?foo=1", b1().SetPath("/info//about").String())

	// add path
	ts.Equal("http:///?foo=1", b1().String())
	ts.Equal("http:///info?foo=1", b1().AddPath("info").String())
	ts.Equal("http:///info/about?foo=1", b1().AddPath("info").AddPath("about").String())
	ts.Equal("http:///info/about?foo=1", b1().AddPath("/info").AddPath("about").String())
	ts.Equal("http:///info/about?foo=1", b1().AddPath("info").AddPath("/about").String())
	ts.Equal("http:///info/about?foo=1", b1().AddPath("/info/").AddPath("/about/").String())
}

// TestQuery tests query modifications
func (ts *URLBuilderTests) TestQuery() {
	// create URL builder
	b1 := func() *URLBuilder {
		return MustParseURL("http:///info")
	}

	// set query
	ts.Equal("http:///info", b1().String())
	ts.Equal("http:///info?foo=1", b1().SetQuery("foo", "1").String())
	ts.Equal("http:///info?foo=2", b1().SetQuery("foo", "1").SetQuery("foo", "2").String())
	ts.Equal("http:///info?foo=1&foo=1", b1().AddQuery("foo", "1").AddQuery("foo", "1").String())
}
