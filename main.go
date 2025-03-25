package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const upcomingURL = "https://www.coasttocoastam.com/shows/upcoming/"

type Show struct {
	Title      string   `json:"title,omitempty"`
	Date       string   `json:"date,omitempty"`
	Host       string   `json:"host,omitempty"`
	Guests     []string `json:"guests,omitempty"`
	FirstHalf  string   `json:"first_half,omitempty"`
	SecondHalf string   `json:"second_half,omitempty"`
}

func main() {
	resp, err := http.Get(upcomingURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	shows := []Show{}
	doc.Find(".feed-cards .coast-feed-item").Each(func(i int, s *goquery.Selection) {
		show := Show{
			Title:      s.Find(".component-heading").Text(),
			Date:       s.Find("time").First().Text(),
			Host:       s.Find(".coast-linked-host .linked-value").Text(),
			FirstHalf:  s.Find(".item-summary p:nth-of-type(1)").Text(),
			SecondHalf: s.Find(".item-summary p:nth-of-type(2)").Text(),
		}

		guests := []string{}
		s.Find(".coast-linked-guests .linked-guest-value a").Each(func(i int, guest *goquery.Selection) {
			guests = append(guests, guest.Text())
		})
		show.Guests = guests

		shows = append(shows, show)
	})

	output, err := json.MarshalIndent(shows, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}
