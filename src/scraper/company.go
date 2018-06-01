package scraper

import (
	"fmt"
)

type EPS_T struct {
	Year  int64
	Value float64
}

type Company struct {
	ID                          int64
	Name                        string
	SupervisorHoldingRatio      float64 // latest
	ForeignInvestorHoldingRatio float64 // latest
	OperatingProfitRatio        float64 // average of four quarter
	EPS                         []EPS_T
	PERatio                     float64 // last year
	ShareCapital                int64
}

func (c *Company) Init(id int64) {
	c.ID = id
	c.EPS = make([]EPS_T, 0)

}

func (c *Company) SetSupervisorHoldingRatio() {
	c.SupervisorHoldingRatio = Config.Parser.ParseSupervisorHoldingRatio(c.ID)
}
func (c *Company) SetForeignInvestorHoldingRatio() {
	c.ForeignInvestorHoldingRatio = Config.Parser.ParseForeignInvestorHoldingRatio(c.ID)
}
func (c *Company) SetOperatingProfitRatio() {
	c.OperatingProfitRatio = Config.Parser.ParseOperatingProfitRatio(c.ID)
}
func (c *Company) SetEPS() {
	c.EPS = Config.Parser.ParseEPS(c.ID)
}

func (c Company) String() string {
	return fmt.Sprintf("ID:%v\n", c.ID) +
		fmt.Sprintf("Supervisor Shareholding Ratio:%v%%\n", c.SupervisorHoldingRatio) +
		fmt.Sprintf("Foreign Investor Ratio:%v%%\n", c.ForeignInvestorHoldingRatio) +
		fmt.Sprintf("(Operating Profit)/(Profit) Ratio:%v%%\n", c.OperatingProfitRatio) +
		fmt.Sprintf("EPS:%v", c.EPS) +
		fmt.Sprintf("Share Capital:%v", c.ShareCapital) +
		fmt.Sprintf("P/E Ratio (last year):%v", c.PERatio)
}
