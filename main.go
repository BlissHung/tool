package main

import (
	"company"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
	"system"
	"time"
	"user"
)

type Company = company.Company

const (
	URL_CHIP_FORMAT = "http://fund.bot.com.tw/Z/ZC/ZCJ/ZCJ_%v.DJHTM"
	//URL_REVENUE_FORMAT = "http://fund.bot.com.tw/z/Z/ZC/ZCH/ZCH_%v.DJHTM"
	URL_PROFIT_FORMAT = "http://fund.bot.com.tw/z//Z/ZC/ZCE/ZCE_%v.DJHTM"
	URL_EPS_FORMAT    = "http://fund.bot.com.tw/z/zc/zcdj_%v.djhtm"
)

func getURL(fromPattern string, id int64) string {
	return fmt.Sprintf(fromPattern, id)
}

var FakeUserAgent = "User-Agent: Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"

func setProfit(c *Company) {

	sumExternalProfit := 0.0
	sumInternalProfit := 0.0

	externalProfiltQueryString := "#oMainTable > tbody > tr:nth-child(%v) > td:nth-child(8)"
	internalProfiltQueryString := "#oMainTable > tbody > tr:nth-child(%v) > td:nth-child(6)"
	profitCollector := colly.NewCollector()
	profitCollector.UserAgent = FakeUserAgent

	getCallOnHTML := func(sum *float64) func(e *colly.HTMLElement) {
		CallOnHTML := func(e *colly.HTMLElement) {
			q := e.DOM
			data, _ := strconv.ParseFloat(strings.Trim(q.Text(), "%"), 64)
			*sum += data
		}
		return CallOnHTML
	}

	FuncExternalProfit := func(index int) {
		format := fmt.Sprintf(externalProfiltQueryString, index)
		profitCollector.OnHTML(format, getCallOnHTML(&sumExternalProfit))
	}
	FuncInternalProfit := func(index int) {
		format := fmt.Sprintf(internalProfiltQueryString, index)
		profitCollector.OnHTML(format, getCallOnHTML(&sumInternalProfit))
	}

	FuncExternalProfit(4)
	FuncExternalProfit(5)
	FuncExternalProfit(6)
	FuncExternalProfit(7)
	FuncInternalProfit(4)
	FuncInternalProfit(5)
	FuncInternalProfit(6)
	FuncInternalProfit(7)
	URL_PROFIT := getURL(URL_PROFIT_FORMAT, c.ID)
	profitCollector.Visit(URL_PROFIT)

	c.OperatingProfitRatio = sumInternalProfit / (sumExternalProfit + sumInternalProfit)
}

var CompanyNameID = make(map[int64]string)

func BuildCompanyID() {
	URL := "http://www.twse.com.tw/exchangeReport/BWIBBU_d?response=html&date=20180524&selectType=ALL"
	c := colly.NewCollector()
	c.UserAgent = FakeUserAgent

	c.OnHTML("body > div > table > tbody > tr", func(e *colly.HTMLElement) {
		q := e.DOM.Children()
		id, _ := q.Html()
		qq := q.NextFiltered("td")
		name, _ := qq.Html()
		idInt, _ := strconv.ParseInt(id, 10, 32)
		CompanyNameID[idInt] = name
	})

	c.Visit(URL)
}

func setChip(c *Company) {

	chipCollector := colly.NewCollector()
	chipCollector.UserAgent = FakeUserAgent
	chipCollector.OnHTML("table.t01>tbody>tr:nth-child(3)>td:nth-child(4)", func(e *colly.HTMLElement) {
		q := e.DOM
		data, _ := strconv.ParseFloat(strings.Trim(q.Text(), "%"), 64)
		c.BoarderRate = data
	})
	chipCollector.OnHTML("table.t01>tbody>tr:nth-child(4)>td:nth-child(4)", func(e *colly.HTMLElement) {
		q := e.DOM
		data, _ := strconv.ParseFloat(strings.Trim(q.Text(), "%"), 64)
		c.ForeignRate = data
	})
	URL_CHIP := getURL(URL_CHIP_FORMAT, c.ID)
	chipCollector.Visit(URL_CHIP)
}

func setPastEPS(c *Company) bool {
	cc := colly.NewCollector()
	cc.UserAgent = FakeUserAgent
	cc.OnHTML("#oMainTable > tbody", func(e *colly.HTMLElement) {
		q := e.DOM
		q = q.Find("tr").Each(func(idx int, s *goquery.Selection) {
			if idx >= 2 {
				year_s := s.Find("td").Eq(0).Text()
				year, _ := strconv.ParseInt(year_s, 10, 32)
				eps_s := s.Find("td").Eq(7).Text()
				eps, _ := strconv.ParseFloat(strings.Trim(eps_s, "%"), 64)

				ele := company.EPS_T{year, eps}
				c.EPS = append(c.EPS, ele)
			}
		})
	})
	URL := getURL(URL_EPS_FORMAT, c.ID)
	cc.Visit(URL)
	return true

}

func runCompany(id int64) {
	// sleep one second before making a request.
	time.Sleep(time.Second)
	company := new(Company)
	company.Init(id)
	if !setPastEPS(company) {
		return
	}
	setChip(company)
	setProfit(company)
	fmt.Println(company)
}

func runAllCompany() {

	// fliter black list
	for id, value := range CompanyNameID {
		if _, ok := user.BlackList[id]; !ok {
			fmt.Println(id, value)
			runCompany(id)
		}
	}
}

func main() {
	BuildCompanyID()

	system.Config.Init()

	if *system.Config.RunAll {
		runAllCompany()
	} else {
		runCompany(*system.Config.RunSingleID)
	}
}
