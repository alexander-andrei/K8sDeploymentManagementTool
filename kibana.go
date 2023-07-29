package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type QueryResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source struct {
				Message string `json:"message"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func KibanaErrorRate() (errorRate float64, err error) {
	kibanaEndpoint := "http://kibana.host.net"
	elasticIndex := "elastic-index-right-here"

	query := `{
		"query": {
			"bool": {
				"must": [
					{
						"match": {
							"level": "error"
						}
					},
					{
						"range": {
							"@timestamp": {
								"gte": "2023-07-01T00:00:00.000Z",
								"lte": "2023-07-28T23:59:59.999Z"
							}
						}
					}
				]
			}
		}
	}`

	resp, err := http.Post(kibanaEndpoint+"/api/console/proxy?path=/"+elasticIndex+"/_search", "application/json", bytes.NewBuffer([]byte(query)))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 400 {
		panic("Cannot connect to kibana")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response QueryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	totalLogs := response.Hits.Total.Value
	errorLogs := len(response.Hits.Hits)

	if totalLogs > 0 {
		errorRate = float64(errorLogs) / float64(totalLogs) * 100.0
		fmt.Printf("Error Rate: %.2f%%\n", errorRate)

	} else {
		fmt.Println("No logs found in the specified time range.")
	}

	return errorRate, err
}
