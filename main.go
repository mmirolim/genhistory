package main

import (
	"fmt"
	"log"
	"time"
)

var (
	// population number
	pN = 10000
	// epoch start
	epochStart = "2015-10-11"
	epS, _     = time.Parse("2006-01-02", epochStart)
	// current day for live a day func
	currentDate = epS
	// email domain group, profile email providers
	edg = []string{"yahoo", "aol", "hotmail", "google"}
	// registration rate
	rr = 0.005
	// unsubscribe rate
	ur = 0.002
	// click rate
	cr = 0.0338
	// there are two types of click weights
	adclick  = 1.0
	conclick = 0.2
	// number of history days to generate
	days = 42 // 6 weeks
	// last user id
	lid int64 = 0
)

func init() {
	log.Println("connect to mongodb")
	err := mongo("localhost:27017", "click_history")
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var err error
	fmt.Println("start generation")
	// on day 0 prepare population
	// generate profiles
	err = genProfiles(10)
	if err != nil {
		log.Fatalln("profiles generation failed", err)
	}
	// live a day from day first
	for d := 1; d < days; d++ {
		if err = liveADay(); err != nil {
			log.Println("liveADay encountered an error", err)
			break
		}

	}
}

// genProfiles creates random profiles and save collection
// to mongodb
func genProfiles(n int) error {
	var err error
	for i := 0; i < n; i++ {
		log.Printf("%+v\n", NewProfile(false))
	}
	return err
}

// one day from site audience
func liveADay() error {
	var err error

	return err
}
