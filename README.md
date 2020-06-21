# go-epp-fi
Golang library and client for interacting with Finnish communication authority's (slightly non-standard) EPP API for FI domains.

## Work In Progress

This repository is a work in progress, with things expected to be changed, moved, added and removed before version 1.0 is released.

Currently done:
- FI EPP extensions (login, logout, balance checking, polling & acking messages)
- Contacts (check, create, read, update, delete)
- Domains (check, create, read, update, delete, renew, transfer)
- Test base for easily test new functions work as they should.

I'm planning to proceed in the following order:
- Host objects (create, read, update, delete)
- FI EPP specialities (transfer lock, DNSSec)
- Small client for controlling contacts & domains.

## Structure

Types for EPP objects can be found under pkg/epp.
Client functionality (that utilizes EPP objects) is available under pkg/registry.

## Tests

Tests for client functionality can be run after local certificates have been created.
Certificate creation has been scripted in Makefile, and creation happens by running `make create-test-certs`.
After this, all tests can be run with the command `make test`.

OpenSSL is required for certificate creation, but tests themselves won't need it.

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
