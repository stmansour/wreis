TOP=../..
THISDIR="pic"

pic: *.go config.json
	@touch fail
	@go vet
	golint
	staticcheck
	go build
	@rm -f fail
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f conf*.json log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? ${THISDIR} fail
	@echo "*** CLEAN completed in ${THISDIR} ***"

config.json:
	cp ../../server/config.json .

relink:
	go build
	@echo "*** RELINK completed in ${THISDIR} ***"

test: pic
	@touch fail
	./pic
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
