package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ingler3/fpr/generate"
	"github.com/ingler3/fpr/search"
	"googlemaps.github.io/maps"
	"log"
	"os"
)

var apiKey = os.Getenv("GOOGLE_API_KEY")
var mapboxAccessToken = os.Getenv("MAPBOX_ACCESS_TOKEN")

var (
	optDest     = flag.String("dest", "dist", "Output directory")
	optKeyword  = flag.String("keyword", "pizza", "Keyword to search for")
	optType     = flag.String("type", "restaurant|bakery", "Types to search for (delimited|by|pipe")
	optLocation = flag.String("location", "Portland, OR", "Location")
)

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()

	if apiKey == "" {
		usageAndExit("Please specify GOOGLE_API_KEY")
	} else if mapboxAccessToken == "" {
		usageAndExit("Please specify MAPBOX_ACCESS_TOKEN")
	} else if optLocation == nil {
		usageAndExit("-location is required")
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed creating client: %s", err)
	}

	log.Println("Finding the results...")
	searchOpts := search.Options{
		Address: *optLocation,
		Type:    *optType,
		Keyword: *optKeyword,
		Limit:   100,
		Radius:  10000, // m
	}

	s, err := search.Do(&searchOpts, c)
	if err != nil {
		log.Fatalf("failed searching: %s", err)
	}

	log.Println("Serializing results...")
	serializedResults, err := json.Marshal(s.Places)
	if err != nil {
		log.Fatalf("failed serializing results: %s", err)
	}

	log.Println("Printing up your new map...")
	compactResults := new(bytes.Buffer)
	if err := json.Compact(compactResults, serializedResults); err != nil {
		log.Fatalf("failed compacting JSON results: %s", err)
	}

	pageData, err := generate.LoadPageData("vars.json")
	if err != nil {
		log.Fatalf("failed loading page vars: %s", err)
	}

	pageData.MapboxAccessToken = mapboxAccessToken
	pageData.Lat = s.LatLng.Lat
	pageData.Lng = s.LatLng.Lng
	pageData.Data = compactResults.String()
	if err := generate.Do(dir, pageData); err != nil {
		log.Fatalf("failed building template: %s", err)
	}
	log.Printf("We're done! Find the goods in %s...\n", *optDest)
}
