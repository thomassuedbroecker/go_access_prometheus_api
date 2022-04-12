# Access Prometheus API 

Simply just access a existing Prometheus REST API with Go.

### How to get started?

* Init the module

```sh
go mod init example/access_prom_api
```

* Write a `main.go` file

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Use Prometheus API to get values")
}
```

* Build a `main.go` file

```sh
go build .
```

* Run a `main.go` 

```sh
go run ./main.go
```
