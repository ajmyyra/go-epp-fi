# go-epp-fi
![Test results](https://github.com/ajmyyra/go-epp-fi/workflows/Tests/badge.svg)

Golang client library, Go API bindings and command line client for interacting with Finnish communication authority's (slightly non-standard) EPP API for FI domains.

## Version 1.0

Version 1.0 has API bindings and client library support for almost if not all FI EPP API functions.
Command line client included in 1.0 has support for balance, service messages, and contact & domain management.

Client library and API bindings contains:
- FI EPP extensions (login, logout, balance checking, polling & acking messages)
- Contacts (check, create, read, update, delete)
- Domains (check, create, read, update, delete, renew, transfer)
- Host objects (create, read, update, delete)
- Tests for almost all library functions, actually querying a local test API.
- Some FI EPP specialities (transfer lock)
- DNSSec support

## Version 1.1

Version 1.1 is a dependency upgrade with one feature added for CLI.

- CLI support for nameserver IP glue records for both IPv4 and IPv6 addresses when transferring, registering or updating domains.

## Version 1.2

Some ideas for the next version. Need something else? Add an [issue](https://github.com/ajmyyra/go-epp-fi/issues) or [pull request](https://github.com/ajmyyra/go-epp-fi/pulls)!

- CLI support for XML debugging (i.e. "send this XML file to server, print what comes back")
- CLI DNSSec support
- Better documentation with examples for GoDoc
- Bubbling under: tests for CLI

## Using the command line client

### Installation

Command line client binary for Linux can be downloaded from Github.

```shell script
$ wget https://github.com/ajmyyra/go-epp-fi/releases/download/1.1.0/epp-fi
$ sudo mv epp-fi /usr/local/bin/
$ chmod +x /usr/local/bin/epp-fi
$ epp-fi
Command line client for FI EPP asset management (domains, contacts etc)

Usage:
  epp-fi [command]

Available Commands:
  balance     Show account balance
  contact     Show, create, edit & delete contacts
  domain      Show, create, edit, transfer & delete domains
  help        Help about any command
  login       Login to FI EPP system
  logout      Logout from FI EPP system (OBS: This cuts all sessions from the user)
  msg         Show and acknowledge service messages

Flags:
  -h, --help   help for epp-fi

Use "epp-fi [command] --help" for more information about a command.
```

### Usage example

CLI configuration can either be stored in a file (.fi-epp.yml) or as environment variables.

Following example sets env variables in a shell session before logging in. Production and debug settings have been commented out.

```shell script
$ cat .env
export FI_EPP_CLIENT_KEY="/home/user/path/to/your_client_key.pem"
export FI_EPP_CLIENT_CERT="/home/user/path/to/your_client_certificate.pem"
export FI_EPP_USERNAME="YOUR_USERNAME"
export FI_EPP_PASSWORD="YOUR_PASSWORD"
# export FI_EPP_SERVER="epp.domain.fi" # For production usage
export FI_EPP_SERVER="epptest.ficora.fi" # For testing purposes
# export FLUME='{"level":"DBG", "development":true, "addCaller":true}' # For debugging

$ source .env

$ epp-fi login
Successfully logged in.

$ epp-fi msg show
ID:       2d159ea5-87fc-4fb9-a792-abd3002fd9a7
Count:    86
Time:     2020-06-07 18:54:13.037 +0000 UTC
Message:  Contact deleted
Name:     C574595

$ # Messages (and other objects) can also be fetches as JSON
$ epp-fi msg show -j
{
  "count": 86,
  "id": "2d159ea5-87fc-4fb9-a792-abd3002fd9a7",
  "date": "2020-06-07T18:54:13.037Z",
  "msg": "Contact deleted",
  "name": "C574595"
}

$ epp-fi msg ack 2d159ea5-87fc-4fb9-a792-abd3002fd9a7
Successfully acknowledged message 2d159ea5-87fc-4fb9-a112-abd3002fd9a7, 85 messages remaining in queue.

$ epp-fi logout
Successfully logged out.
```

## Project structure

Types for EPP objects can be found under pkg/epp.
Client functionality (that utilizes EPP objects) is available under pkg/registry.
Command line client (that utilizes the EPP objects and client) is under cmd.

## Tests

Tests for client functionality can be run after local certificates have been created.
Certificate creation has been scripted in Makefile, and creation happens by running `make create-test-certs`.
After this, all tests can be run with the command `make test`.

OpenSSL is required for certificate creation, but tests themselves won't need it.

