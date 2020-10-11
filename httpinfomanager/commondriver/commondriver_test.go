package commondriver

import "testing"
import "github.com/herb-go/herbmodules/httpinfomanager"

func TestDriver(t *testing.T) {
	formatterdrivers := []string{
		"toupper",
		"tolower",
		"trim",
		"integer",
		"match",
		"find",
		"split",
		"hired",
	}
	for _, v := range formatterdrivers {
		d, err := httpinfomanager.GetFormatterFactory(v)
		if d == nil || err != nil {
			t.Fatal(d, err)
		}
	}
	extractordrivers := []string{
		"header",
		"query",
		"form",
		"router",
		"fixed",
		"cookie",
		"ip",
		"method",
		"path",
		"host",
		"user",
		"password",
		"identifier",
		"hired",
	}
	for _, v := range extractordrivers {
		d, err := httpinfomanager.GetExtractorFactory(v)
		if d == nil || err != nil {
			t.Fatal(d, err)
		}
	}
}
