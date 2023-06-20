TOP=../..
THISDIR="lib"

lib: *.go config.json
	@touch fail
	@go vet
	golint
	staticcheck
	go build
	# go test
	go install
	@rm -f fail
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? ${THISDIR} fail conf*.json
	@echo "*** CLEAN completed in ${THISDIR} ***"

relink:
	go build
	@echo "*** RELINK completed in ${THISDIR} ***"

config.json:
	@cp ${TOP}/server/config.json config.json

test:
	@touch fail
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
