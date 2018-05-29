package system

import "flag"

type globalConfig struct {
	RunAll      *bool
	RunSingleID *int64
}

var Config globalConfig

func (config *globalConfig) Init() {

	config.RunSingleID = flag.Int64("id", 2498, "Specify a company ID")
	config.RunAll = flag.Bool("runAll", false, "Run all company IDs")
	flag.Parse()

}
