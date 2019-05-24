package main

import (
	"fmt"
	"honnef.co/go/js/dom"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

type spec struct {
	URL, Notes       string
	MaxTemp, MinTemp float64
}

var tags = map[string]spec{
	"Ambient": {
		URL:     "https://my.wirelesstag.net/ethSharedFrame.aspx?pic=1&hide_name=0&hide_temp=0&hide_rh=0&hide_updated=0&hide_signal=1&hide_battery=0&hide_battery_volt=1&hide_motion=1&uuids=784d6a45-f662-416e-9707-97241cc69b2d",
		Notes:   "Ambient temperature",
		MinTemp: 10,
		MaxTemp: 30,
	},
	"Projector": {
		URL:     "https://my.wirelesstag.net/ethSharedFrame.aspx?pic=1&hide_name=0&hide_temp=0&hide_rh=0&hide_updated=0&hide_signal=1&hide_battery=0&hide_battery_volt=1&hide_motion=1&uuids=291ba7e4-9606-4cc7-bec0-ebf6e023628e",
		Notes:   "Keep me cool please!",
		MinTemp: 0,
		MaxTemp: 30,
	},
	"Food esky": {
		URL:     "https://my.wirelesstag.net/ethSharedFrame.aspx?pic=1&hide_name=0&hide_temp=0&hide_rh=0&hide_updated=0&hide_signal=1&hide_battery=0&hide_battery_volt=1&hide_motion=1&uuids=0b8a8a95-968b-4918-a8d2-147c9684b80c",
		Notes:   "Below 5 C",
		MinTemp: 0,
		MaxTemp: 5,
	},
	"Meat store": {
		URL:     "https://my.wirelesstag.net/ethSharedFrame.aspx?pic=1&hide_name=0&hide_temp=0&hide_rh=0&hide_updated=0&hide_signal=1&hide_battery=0&hide_battery_volt=1&hide_motion=1&uuids=bef248e5-29d4-4443-b795-8ac2d53ec98e",
		Notes:   "Below 5 C",
		MinTemp: 10,
		MaxTemp: 20,
	},
	"Cask store": {
		URL:     "https://my.wirelesstag.net/ethSharedFrame.aspx?pic=1&hide_name=0&hide_temp=0&hide_rh=0&hide_updated=0&hide_signal=1&hide_battery=0&hide_battery_volt=1&hide_motion=1&uuids=a8564fa2-ee17-4895-81d3-754c0d11bc9f",
		Notes:   "8-12 C",
		MinTemp: 10,
		MaxTemp: 20,
	},
	"Bar casks": {
		URL:     "https://my.wirelesstag.net/ethSharedFrame.aspx?pic=1&hide_name=0&hide_temp=0&hide_rh=0&hide_updated=0&hide_signal=1&hide_battery=0&hide_battery_volt=1&hide_motion=1&uuids=59c37fd7-ec15-4dce-9557-6d4e97f67005",
		Notes:   "12C. Range 10-14 C",
		MinTemp: 10,
		MaxTemp: 5,
	},
	"Bar key kegs": {
		URL:     "https://my.wirelesstag.net/ethSharedFrame.aspx?pic=1&hide_name=0&hide_temp=0&hide_rh=0&hide_updated=0&hide_signal=1&hide_battery=0&hide_battery_volt=1&hide_motion=1&uuids=863b88f1-b318-456c-9577-ff772be2607c",
		Notes:   "8C. Range 6-10 C",
		MinTemp: 10,
		MaxTemp: 20,
	},
}

func main() {

	doc := dom.GetWindow().Document()
	//holder := doc.GetElementsByTagName("body")[0]
	holder := doc.GetElementByID("p1")

	// Add a table
	//table := doc.CreateElement("table")
	//holder.AppendChild(table)

	// Add a header row
	row := doc.CreateElement("tr")
	//table.AppendChild(row)
	holder.AppendChild(row)
	nameCell := doc.CreateElement("td")
	tempCell := doc.CreateElement("td")
	row.AppendChild(nameCell)
	row.AppendChild(tempCell)
	//nameCell.SetInnerHTML("Name")
	//tempCell.SetInnerHTML("Temperature")

	for name, data := range tags {

		// Add a row to the table
		row := doc.CreateElement("tr")
		//table.AppendChild(row)
		holder.AppendChild(row)
		nameCell := doc.CreateElement("td")
		tempCell := doc.CreateElement("td")
		notesCell := doc.CreateElement("td")
		notesCell.SetAttribute("style", "color: Gray;")

		row.AppendChild(nameCell)

		row.AppendChild(tempCell)
		row.AppendChild(notesCell)
		nameCell.SetInnerHTML(name)

		tempCell.SetInnerHTML("loading...")
		notesCell.SetInnerHTML(data.Notes)

		go getTag(name, data, tempCell)
	}
}

var rxp = regexp.MustCompile(`(?:F\/)(.*)(?:\&deg;C)`)

func getTag(name string, data spec, cell dom.Element) {
	resp, err := http.Get(data.URL)
	if err != nil {
		cell.SetInnerHTML(err.Error())
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cell.SetInnerHTML(err.Error())
		return
	}

	matches := rxp.FindSubmatch(b)
	if len(matches) < 2 {
		cell.SetInnerHTML("can't find temp")
		return
	}

	temp, err := strconv.ParseFloat(string(matches[1]), 64)
	if err != nil {
		cell.SetInnerHTML(err.Error())
		return
	}

	cell.SetInnerHTML(fmt.Sprintf("%.1fC", temp))

	if temp < data.MinTemp {
		//cell.SetAttribute("style", "border: 2px blue solid;")
		cell.SetAttribute("style", "color: blue;")
	} else if temp > data.MaxTemp {
		//cell.SetAttribute("style", "border: 2px red solid;")
		cell.SetAttribute("style", "color: red;")
	}

	return
}
