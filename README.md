# go-epp-fi
Golang library and client for interacting with Finnish communication authority's (slightly non-standard) EPP API for FI domains.

## Work In Progress

This repository is a work in progress, with things expected to be changed, moved, added and removed before version 1.0 is released.

Currently done:
- FI EPP extensions (login, logout, balance checking, polling & acking messages)
- Contacts (check, create, read, update, delete)
- Domains (check, create, read, update)

I'm planning to proceed in the following order:
- Domains (delete, renew, transfer)
- Host objects (create, read, update, delete)
- FI EPP specialities (transfer lock, DNSSec)

## Raw example usage without any features

```
import (
	"fmt"
	"github.com/ajmyyra/go-epp-fi/pkg/registry"
	"io/ioutil"
)

func main() {
	clientKey, err := ioutil.ReadFile("/path/to/your/certs/privkey.pem")
	if err != nil {
		panic(err)
	}
	clientCert, err := ioutil.ReadFile("/path/to/your/certs/cert.pem")
    if err != nil {
		panic(err)
	}

	client, err := registry.NewRegistryClient("foo","bar","epptest.ficora.fi", 700, clientKey, clientCert)
	if err != nil {
		panic(err)
	}

	if err = client.Connect(); err != nil {
		panic(err)
	}
	fmt.Printf("Connected successfully, server time is now %s\n", client.Greeting.SvDate)

	greeting, err := client.Hello()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", greeting)

    if err = client.Login(); err != nil {
        panic(err)
    }
    fmt.Println("Logged in successfully.")

    if err = client.Logout(); err != nil {
        panic(err)
    }
    fmt.Println("Logged out successfully.")
}
```
