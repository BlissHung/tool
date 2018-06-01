package scraper

import (
	"flag"
)

type globalConfig struct {
	RunAll      *bool
	RunSingleID *int64
	Parser      InfoParser
}

var Config globalConfig

func (config *globalConfig) Init() {

	config.RunSingleID = flag.Int64("id", 2498, "Specify a company ID")
	config.RunAll = flag.Bool("runAll", false, "Run all company IDs")

	p := new(TaiwanBankParser)
	config.Parser = p
	flag.Parse()

}
