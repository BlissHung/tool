package main

import (
	"fmt"
	"scraper"
	"userConfig"
)

type Company = scraper.Company

func runCompany(id int64) {
	company := new(Company)
	company.Init(id)
	company.SetSupervisorHoldingRatio()
	company.SetForeignInvestorHoldingRatio()
	company.SetOperatingProfitRatio()
	company.SetEPS()
	fmt.Println(company)
}

func runAllCompany() {

	// fliter black list
	for id, value := range scraper.Config.Parser.ParseCompanyList() {
		if _, ok := userConfig.BlackList[id]; !ok {
			fmt.Println(id, value)
			runCompany(id)
		}
	}
}

func main() {

	scraper.Config.Init()

	if *scraper.Config.RunAll {
		runAllCompany()
	} else {
		runCompany(*scraper.Config.RunSingleID)
	}
}
