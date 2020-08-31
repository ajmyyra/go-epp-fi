# go-epp-fi
![Test results](https://github.com/ajmyyra/go-epp-fi/workflows/Tests/badge.svg)

Golang client library, Go API bindings and command line client for interacting with Finnish communication authority's (slightly non-standard) EPP API for FI domains.

## Version 1.0

Version 1.0 has API bindings and client library support for all FI EPP API functions.
Command line client included in 1.0 has support for balance, service messages, and contact & domain management.

Client library and API bindings support:
- FI EPP extensions (login, logout, balance checking, polling & acking messages)
- Contacts (check, create, read, update, delete)
- Domains (check, create, read, update, delete, renew, transfer)
- Host objects (create, read, update, delete)
- Base for all API tests and tests for FI EPP extensions, contacts & domains
- Some FI EPP specialities (transfer lock)
- DNSSec support
- Command line client (epp-fi) for controlling contacts & domains, checking service messages and checking balance


## Version 1.1

Some ideas for the next version. Need something else? Add an [issue](https://github.com/ajmyyra/go-epp-fi/issues) or [pull request](https://github.com/ajmyyra/go-epp-fi/pulls)!

- CLI support for XML debugging (i.e. "send this XML file to server, print what comes back")
- CLI DNSSec support
- Better documentation with examples for GoDoc
- Bubbling under: tests for CLI

## Project structure

Types for EPP objects can be found under pkg/epp.
Client functionality (that utilizes EPP objects) is available under pkg/registry.
Command line client (that utilizes the EPP objects and client) is under cmd.

## Tests

Tests for client functionality can be run after local certificates have been created.
Certificate creation has been scripted in Makefile, and creation happens by running `make create-test-certs`.
After this, all tests can be run with the command `make test`.

OpenSSL is required for certificate creation, but tests themselves won't need it.

