package commondriver

import (
	"github.com/herb-go/herbmodules/httpinfomanager/extractors/commonextractor"
	"github.com/herb-go/herbmodules/httpinfomanager/extractors/workerextractor"
	"github.com/herb-go/herbmodules/httpinfomanager/formatters/commonformatter"
	"github.com/herb-go/herbmodules/httpinfomanager/formatters/workerformatter"
)

func RegsiterFactories() {
	commonextractor.RegsiterFactories()
	workerextractor.RegsiterFactories()
	commonformatter.RegisterFactories()
	workerformatter.RegisterFactories()

}
func init() {
	RegsiterFactories()
}
