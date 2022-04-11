package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// ***************************
// Struct build on'example-prometheus-query-result.json'
// for later unmarshal json
// -> start

type QueryResult struct {
	Status string
	Data   Data
}

type Data struct {
	ResultType string
	Result     []Result
}

type Result struct {
	Metric Metric
}

type Metric struct {
	__Name__  string
	Container string
	Endpoint  string
	Instance  string
	Job       string
	Namespace string
	Pod       string
	Service   string
	Value     []string
	// Value     []interface{}
}

// -> end
// ***************************

func main() {
	fmt.Println("Use Prometheus API to get values")

	goobers_query := "goobers_total"
	prom_api_call := "api/v1/query?query="
	PROMETHEUS_URL := ""
	// request_url := os.Getenv("PROMETHEUS_URL") + "/" + prom_api_call + goobers_query
	request_url := PROMETHEUS_URL + "/" + prom_api_call + goobers_query
	fmt.Println("Request url : " + request_url)

	// 1. Create HTTP request
	req, err := http.NewRequest("GET", request_url, nil)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	} else {
		fmt.Println("Request Host: " + req.Host)
	}

	// 2. Define header
	req.Header.Set("Accept", "application/json")

	// 3. Create client
	client := http.Client{}

	// 4. Invoke HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	fmt.Printf("resp.Header: %v\n", resp.Header)
	fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)

	// 5. Verify the request status
	if resp.StatusCode == http.StatusOK {

		// 6. Get only body from response
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		bodyString := string(bodyBytes)
		fmt.Printf("bodyString: %v\n", bodyString)

		// 7. Convert body to json content
		var dat QueryResult
		if err := json.Unmarshal(bodyBytes, &dat); err != nil {
			panic(err)
		}

		fmt.Printf("bodyString: %v\n", dat.Data.ResultType)
		fmt.Printf("bodyString: %v\n", dat.Data.Result[0].Metric.__Name__)
		/*
			if counter_float, ok := dat.Data.Result[0].Metric.Value[0].(float32); ok {
				// act on float
				fmt.Printf("counter_float: %f\n", counter_float)
			} else {
				// not float
				fmt.Printf("Error counter_float")
			}

			if counter_string, ok := dat.Data.Result[0].Metric.Value[1].(string); ok {
				// act on str
				fmt.Printf("counter_string: %v\n", counter_string)
			} else {
				// not string
				fmt.Printf("Error counter_float")
			}
		*/
	}

	fmt.Println("**********DONE********")

}
