package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-martini/martini"
	"github.com/hackedu/backend/v1/helper"
)

type Marker struct {
	Location Location `json:"marker"`
}

type Driver struct {
	ID      int64      `json:"id,string"`
	History []Location `json:"history"`
}

type Location struct {
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"lat"`
	Longitude float64   `json:"lng"`
}

type LocationResponse struct {
	Drivers []DriverResponse `json:"drivers"`
}

type DriverResponse struct {
	ID        int64   `json:"id,string"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

const (
	apiBase     = "https://api.lyft.com"
	apiLocation = apiBase + "/users/%d/location"
)

var DriverRecords map[int64]*Driver

func main() {
	go webserver()

	comm := make(chan map[int64]*DriverResponse)

	go func(c chan map[int64]*DriverResponse) {
		tick := time.NewTicker(5 * time.Second)
		for _ = range tick.C {
			data, err := poll()
			if err != nil {
				log.Fatal(err)
			}
			comm <- data
		}
	}(comm)

	DriverRecords = make(map[int64]*Driver)

	for {
		select {
		case driverResults := <-comm:
			for key, driver := range driverResults {
				if d, ok := DriverRecords[key]; ok {
					location := d.History[len(d.History)-1]

					if location.Latitude != driver.Latitude ||
						location.Longitude != driver.Longitude {
						d.History = append(d.History, Location{
							Timestamp: time.Now(),
							Latitude:  driver.Latitude,
							Longitude: driver.Longitude,
						})
					}
				} else {
					d := Driver{
						ID: driver.ID,
						History: []Location{
							Location{
								Timestamp: time.Now(),
								Latitude:  driver.Latitude,
								Longitude: driver.Longitude,
							},
						},
					}

					DriverRecords[key] = &d
				}
			}
			log.Println(DriverRecords)
		default:
		}
	}
}

func webserver() {
	m := martini.Classic()

	m.Use(martini.Static("public"))
	m.Get("/api/drivers", func(w http.ResponseWriter) (int, string) {
		drivers := make([]Driver, len(DriverRecords))
		i := 0
		for _, driver := range DriverRecords {
			drivers[i] = *driver
			i++
		}

		json, err := json.Marshal(drivers)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Internal server error"
		}

		w.Header().Set("Content-Type", "application/json")
		return http.StatusOK, string(json)
	})

	m.Run()
}

func poll() (map[int64]*DriverResponse, error) {
	markers := []Marker{
		Marker{
			Location: Location{
				Latitude:  34.07030187403868,
				Longitude: -118.44679702073337,
			},
		},
		Marker{
			Location: Location{
				Latitude:  33.9251508,
				Longitude: -118.4083849,
			},
		},
		Marker{
			Location: Location{
				Latitude:  33.9456383,
				Longitude: -118.393934,
			},
		},
		Marker{
			Location: Location{
				Latitude:  34.0201391,
				Longitude: -118.4985162,
			},
		},
		Marker{
			Location: Location{
				Latitude:  34.0195053,
				Longitude: -118.2875214,
			},
		},
	}

	response := make(chan LocationResponse)
	log.Println("Polling...")

	for _, marker := range markers {
		go func(ch chan LocationResponse, marker Marker) {
			// TOOD: Make own method
			client := new(http.Client)
			markerJson, err := json.Marshal(marker)
			if err != nil {
				log.Println(err)
			}

			markerReader := bytes.NewReader(markerJson)

			userId, err := strconv.Atoi(helper.GetConfig("USER_ID"))
			if err != nil {
				log.Println(err)
			}
			url := fmt.Sprintf(apiLocation, userId)
			req, err := http.NewRequest("PUT", url, markerReader)
			if err != nil {
				log.Println(err)
			}

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "fbAccessToken "+helper.GetConfig("FACEBOOK_ACCESS_TOKEN"))

			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}

			var locationResponse LocationResponse
			err = json.Unmarshal(body, &locationResponse)
			if err != nil {
				log.Println(err)
			}
			ch <- locationResponse
		}(response, marker)
	}

	count := len(markers)
	drivers := make(map[int64]*DriverResponse)
	for {
		select {
		case resp := <-response:
			count--

			for _, driver := range resp.Drivers {
				if _, ok := drivers[driver.ID]; !ok {
					drivers[driver.ID] = &DriverResponse{driver.ID, driver.Latitude, driver.Longitude}
				}
			}
		default:
			if count <= 0 {
				log.Println("Poll complete")
				return drivers, nil
			}
		}
	}
}
