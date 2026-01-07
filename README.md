# teller-go

Go client for [teller.io](https://teller.io) API by [maxint.com](https://maxint.com/)

## Installation

```bash
go get github.com/maxint-app/teller-go
```

## Usage

Go to https://teller.io/dashboard > Certificates > Create a new certificate.

Save certificate files `certificate.pem` and `private_key.pem` in project root.

Add `.pem` to `.gitignore`, to tell Git to ignore certificate files when you make a commit.

Use in your project:

```go
package main

import (
	"fmt"
	"log"

	"github.com/maxint-app/teller-go"
)

func main() {
	client, err := teller.NewClient(
		"./certificate.pem",
		"./private_key.pem",
		nil, // optional access token
	)
	if err != nil {
		log.Fatal(err)
	}

	// optional: pass access token in individual requests
	accessToken := "your_access_token_from_teller_connect"
	identity, err := client.Identity.Get(&teller.TellerOptionsBase{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", identity)
}
```

> Follow the teller.io [docs](https://teller.io/docs/api) for more information.

## License

MIT
