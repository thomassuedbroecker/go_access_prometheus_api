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
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Metric struct {
	Name      string `json:"__name__"`
	Container string `json:"container"`
	Endpoint  string `json:"endpoint"`
	Instance  string `json:"instance"`
	Job       string `json:"job"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	Service   string `json:"service"`
}

type Result struct {
	Metric Metric        `json:"metric"`
	Value  []interface{} `json:"value"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

// -> end
// ***************************

func main() {
	fmt.Println("*******************************")
	fmt.Println("-> Use Prometheus API to get counter values")

	// 0. Define Prometheus REST API request
	goobers_query := "goobers_total"
	prom_api_call := "api/v1/query?query="
	// Insert URL with environment variable
	// export PROMETHEUS_URL=http://HOST:PORT
	request_url := os.Getenv("PROMETHEUS_URL") + "/" + prom_api_call + goobers_query

	// Insert URL inside the script
	// PROMETHEUS_URL := ""
	// request_url := PROMETHEUS_URL + "/" + prom_api_call + goobers_query

	fmt.Println("Prometheus REST API request: " + request_url)

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

	// 5. Provide response information
	fmt.Printf("resp.Header: %v\n", resp.Header)
	fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)

	// 6. Verify the request status
	if resp.StatusCode == http.StatusOK {

		// 7. Get only body from response
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		// Optional inspect JSON
		// inspectJSON(bodyBytes)

		bodyString := string(bodyBytes)

		fmt.Printf("-> Response bodyString: %v\n", bodyString)

		// 8. Convert body to json content
		var dat QueryResult
		if err := json.Unmarshal(bodyBytes, &dat); err != nil {
			panic(err)
		}

		// 9. Show values from the response json
		fmt.Printf("-> dat.Data.ResultType: %v\n", dat.Data.ResultType)
		fmt.Printf("-> dat.Data.Result[0].Metric.Name: %v\n", dat.Data.Result[0].Metric.Name)

		// 10. Decode interface datatypes
		if counter_float, ok := dat.Data.Result[0].Value[0].(float64); ok {
			// act on float
			fmt.Printf("-> counter_float: %f\n", counter_float)
		} else {
			// not float
			fmt.Printf("-> Error counter_float")
		}

		if counter_string, ok := dat.Data.Result[0].Value[1].(string); ok {
			// act on str
			fmt.Printf("-> counter_string: %v\n", counter_string)
		} else {
			// not string
			fmt.Printf("Error counter_float")
		}
	}

	fmt.Println("**********DONE********")
}

func inspectJSON(bodyBytes []byte) {

	var f interface{}

	err := json.Unmarshal(bodyBytes, &f)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float32:
			fmt.Println(k, "is float32", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []string:
			fmt.Println(k, "is []string", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case Result:
			fmt.Println(k, "is Result", vv)
		case interface{}:
			fmt.Println(k, "is interface{}", vv)
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}
