package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

type InfoParser interface {
	ParseCompanyList() map[int64]string
	ParseSupervisorHoldingRatio(int64) float64
	ParseForeignInvestorHoldingRatio(int64) float64
	ParseOperatingProfitRatio(int64) float64
	ParseEPS(int64) []EPS_T
}

func getURL(fromPattern string, id int64) string {
	return fmt.Sprintf(fromPattern, id)
}

const (
	URL_CHIP_FORMAT    = "http://fund.bot.com.tw/Z/ZC/ZCJ/ZCJ_%v.DJHTM"
	URL_REVENUE_FORMAT = "http://fund.bot.com.tw/z/Z/ZC/ZCH/ZCH_%v.DJHTM"
	URL_PROFIT_FORMAT  = "http://fund.bot.com.tw/z//Z/ZC/ZCE/ZCE_%v.DJHTM"
	URL_EPS_FORMAT     = "http://fund.bot.com.tw/z/zc/zcdj_%v.djhtm"
)

var FakeUserAgent = "User-Agent: Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"

type TaiwanBankParser struct{}

func (p *TaiwanBankParser) ParseCompanyList() map[int64]string {

	Ret := make(map[int64]string)
	URL := "http://www.twse.com.tw/exchangeReport/BWIBBU_d?response=html&date=20180524&selectType=ALL"
	c := colly.NewCollector()
	c.UserAgent = FakeUserAgent

	c.OnHTML("body > div > table > tbody > tr", func(e *colly.HTMLElement) {
		q := e.DOM.Children()
		id, _ := q.Html()
		qq := q.NextFiltered("td")
		name, _ := qq.Html()
		idInt, _ := strconv.ParseInt(id, 10, 32)
		Ret[idInt] = name
	})

	c.Visit(URL)
	return Ret
}

func (p *TaiwanBankParser) ParseForeignInvestorHoldingRatio(id int64) float64 {
	chipCollector := colly.NewCollector()
	chipCollector.UserAgent = FakeUserAgent

	var Ret float64 = 0.0
	//	chipCollector.OnHTML("table.t01>tbody>tr:nth-child(4)>td:nth-child(4)", func(e *colly.HTMLElement) {
	chipCollector.OnHTML("table.t01>tbody>tr:nth-child(4)", func(e *colly.HTMLElement) {
		q := e.DOM
		Ret, _ = strconv.ParseFloat(strings.Trim(q.Text(), "%"), 64)
	})
	URL_CHIP := getURL(URL_CHIP_FORMAT, id)
	chipCollector.Visit(URL_CHIP)

	return Ret
}
func (p *TaiwanBankParser) ParseSupervisorHoldingRatio(id int64) float64 {

	chipCollector := colly.NewCollector()
	chipCollector.UserAgent = FakeUserAgent

	var Ret float64 = 0.0
	chipCollector.OnHTML("table.t01>tbody>tr:nth-child(3)>td:nth-child(4)", func(e *colly.HTMLElement) {
		q := e.DOM
		Ret, _ = strconv.ParseFloat(strings.Trim(q.Text(), "%"), 64)
	})
	URL_CHIP := getURL(URL_CHIP_FORMAT, id)
	chipCollector.Visit(URL_CHIP)

	return Ret
}
func (p *TaiwanBankParser) ParseEPS(id int64) []EPS_T {

	Ret := make([]EPS_T, 0)
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

				ele := EPS_T{year, eps}
				Ret = append(Ret, ele)
			}
		})
	})
	URL := getURL(URL_EPS_FORMAT, id)
	cc.Visit(URL)
	return Ret
}

func (p *TaiwanBankParser) ParseOperatingProfitRatio(id int64) float64 {
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
	URL_PROFIT := getURL(URL_PROFIT_FORMAT, id)
	profitCollector.Visit(URL_PROFIT)

	return sumInternalProfit / (sumExternalProfit + sumInternalProfit)
}
