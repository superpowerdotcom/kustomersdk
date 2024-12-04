# kustomersdk

[![Main](https://github.com/superpowerdotcom/kustomersdk/actions/workflows/main.yml/badge.svg)](https://github.com/superpowerdotcom/kustomersdk/actions/workflows/main.yml)
[![PR](https://github.com/superpowerdotcom/kustomersdk/actions/workflows/pr.yml/badge.svg)](https://github.com/superpowerdotcom/kustomersdk/actions/workflows/pr.yml)


Golang SDK for Kustomer API



## Usage

```go
package main

import (
	"context"
	"log"
	
	"github.com/superpowerdotcom/kustomersdk"
)

func main() {
	client, err := kustomersdk.New(&kustomersdk.Options{
		APIKey: "your-api-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get a customer
	customerRecord, err := client.GetCustomerByExternalID(context.Background(), "postgres-database-uuid")
	if err != nil {
		log.Fatal(err)
	}
	
	println(customerRecord.Data.Attributes.Name) // John Smith Jr.
}
```

## Tests

* Mock responses are stored under `test-assets/responses`
* Run tests with `make test`
```