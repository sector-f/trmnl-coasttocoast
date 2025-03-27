package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	upcomingURL = "https://www.coasttocoastam.com/shows/upcoming/"
	timeFormat  = "Monday, January 2, 2006"
)

type Output struct {
	Shows []Show `json:"shows,omitempty"`
}

type Show struct {
	Title      string `json:"title,omitempty"`
	Date       string `json:"date,omitempty"`
	Host       string `json:"host,omitempty"`
	Guests     string `json:"guests,omitempty"`
	FirstHalf  string `json:"first_half,omitempty"`
	SecondHalf string `json:"second_half,omitempty"`
}

func main() {
	log.Printf("Fetching %s\n", upcomingURL)
	resp, err := http.Get(upcomingURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received %v status\n", resp.Status)
	}

	log.Println("Parsing response as HTML")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	out := Output{}
	doc.Find(".feed-cards .coast-feed-item").Each(func(i int, s *goquery.Selection) {
		show := Show{
			Title:      s.Find(".component-heading").Text(),
			Host:       s.Find(".coast-linked-host .linked-value").Text(),
			FirstHalf:  s.Find(".item-summary p:nth-of-type(1)").Text(),
			SecondHalf: s.Find(".item-summary p:nth-of-type(2)").Text(),
		}

		date, err := time.Parse(timeFormat, s.Find("time").First().Text())
		if err == nil {
			show.Date = date.Format(time.DateOnly)
		}

		guests := []string{}
		s.Find(".coast-linked-guests .linked-guest-value a").Each(func(i int, guest *goquery.Selection) {
			guests = append(guests, guest.Text())
		})
		show.Guests = strings.Join(guests, ", ")

		out.Shows = append(out.Shows, show)
	})

	output, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}
