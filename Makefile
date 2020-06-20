.EXPORT_ALL_VARIABLES:
TEMP_TEST_DIR = testtmp
TEMP_CA_DIR = catemp

.DEFAULT: create-test-certs

.PHONY: create-test-certs
create-test-certs:
	mkdir $(TEMP_TEST_DIR)
	mkdir $(TEMP_CA_DIR)
	touch $(TEMP_CA_DIR)/index.txt
	echo 1000 > $(TEMP_CA_DIR)/serial
	openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 -subj "/C=FI/ST=Uusimaa/O=TestOrg/CN=testCA" -keyout $(TEMP_CA_DIR)/rootCA.key -out $(TEMP_CA_DIR)/rootCA.crt
	openssl genrsa -out $(TEMP_TEST_DIR)/testserver.key 2048
	openssl req -new -sha256 -key $(TEMP_TEST_DIR)/testserver.key -subj "/C=FI/ST=Uusimaa/O=TestOrg/CN=testserver" -addext "subjectAltName = IP:127.0.0.1" -out $(TEMP_TEST_DIR)/testserver.csr
	openssl ca -batch -policy policy_anything -keyfile $(TEMP_CA_DIR)/rootCA.key -cert $(TEMP_CA_DIR)/rootCA.crt -out $(TEMP_TEST_DIR)/testserver.crt -config test/openssl.cnf -extensions v3_req -infiles $(TEMP_TEST_DIR)/testserver.csr
	openssl genrsa -out $(TEMP_TEST_DIR)/testclient.key 2048
	openssl req -new -sha256 -key $(TEMP_TEST_DIR)/testclient.key -subj "/C=FI/ST=Uusimaa/O=TestOrg/CN=testclient" -addext "subjectAltName = IP:127.0.0.1" -out $(TEMP_TEST_DIR)/testclient.csr
	openssl ca -batch -policy policy_anything -keyfile $(TEMP_CA_DIR)/rootCA.key -cert $(TEMP_CA_DIR)/rootCA.crt -out $(TEMP_TEST_DIR)/testclient.crt -config test/openssl.cnf -extensions v3_req -infiles $(TEMP_TEST_DIR)/testclient.csr
	mv $(TEMP_CA_DIR)/rootCA.crt $(TEMP_TEST_DIR)/
	rm -rf $(TEMP_CA_DIR)

.PHONY: test
test:
	go test -v ./...

.PHONY: clean-test-ca
clean-test-ca:
	rm -rf $(TEMP_TEST_DIR)
	rm -rf $(TEMP_CA_DIR)

.PHONY: clean
clean: clean-test-ca
