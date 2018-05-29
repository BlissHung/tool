package company

import (
	"fmt"
)

type EPS_T struct {
	Year  int64
	Value float64
}

type Company struct {
	ID                   int64
	Name                 string
	BoarderRate          float64 // latest
	ForeignRate          float64 // latest
	OperatingProfitRatio float64 // average of four quarter
	EPS                  []EPS_T
}

func (c *Company) Init(id int64) {
	c.ID = id
	c.EPS = make([]EPS_T, 0)

}

func (c Company) String() string {
	return fmt.Sprintf("ID:%v\n", c.ID) +
		fmt.Sprintf("Supervisor Shareholding Ratio:%v%%\n", c.BoarderRate) +
		fmt.Sprintf("Foreign Investor Ratio:%v%%\n", c.ForeignRate) +
		fmt.Sprintf("(Operating Profit)/(Profit) Ratio:%v%%\n", c.OperatingProfitRatio) +
		fmt.Sprintf("EPS:%v", c.EPS)
}
