# go-epp-fi
Golang library and client for interacting with Finnish communication authority's (slightly non-standard) EPP API for FI domains.

## Work In Progress

This repository is a work in progress, with things expected to be changed, moved, added and removed before version 1.0 is released.

Currently only connectivity is available. I'm planning to proceed in the following order:
- FI EPP extensions (login, logout, balance checking, polling & acking messages)
- Contacts (create, read, update, delete)
- Domains (create, read, update, delete, renew, transfer)
- Host objects (create, read, update, delete)
- FI EPP specialities (transfer lock, DNSSec)

## Raw example usage without any features

```
clientKey, err := ioutil.ReadFile("/path/to/your/privkey.pem")
if err != nil {
	panic(err)
}
clientCert, err := ioutil.ReadFile("/path/to/your/cert.pem")

client, err := registry.NewRegistryClient("username","passwd","epptest.ficora.fi", 700, clientKey, clientCert)
if err != nil {
	panic(err)
}

if err = client.Connect(); err != nil {
	panic(err)
}

hello := epp.APIHello{
	XMLName: xml.Name{},
	Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
}

helloxml, err := xml.Marshal(hello)
if err != nil {
	panic(err)
}

fmt.Printf("Hello to send: %s\n", string(helloxml))
err = client.Write(helloxml)
if err != nil {
	panic(err)
}

// To give API some time to process our request, not usually necessary.
time.Sleep(1 * time.Second)

reply, err := client.Read()
if err != nil {
	panic(err)
}
```
