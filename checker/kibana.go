package checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s/tool/config"
	"net/http"
	"time"
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
	now := time.Now()
	thirtyMinutesAgo := now.Add(-30 * time.Minute)
	startDate := thirtyMinutesAgo.Format(time.RFC3339Nano)
	endDate := now.Format(time.RFC3339Nano)

	query := fmt.Sprintf(`{
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
								"gte": "%s",
								"lte": "%s"
							}
						}
					}
				]
			}
		}
	}`, startDate, endDate)

	resp, err := http.Post(config.GlobalConfig.Kibana.Endpoint+"/api/console/proxy?path=/"+config.GlobalConfig.Kibana.ElasticIndex+"/_search", "application/json", bytes.NewBuffer([]byte(query)))
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
