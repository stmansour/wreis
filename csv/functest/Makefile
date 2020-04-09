TOP=../..
THISDIR="csv/functest"

functest: *.go
	@touch fail
	@go vet
	@golint
	go build
	@rm -f fail
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? functest fail
	@echo "*** CLEAN completed in ${THISDIR} ***"

relink:
	go build
	@echo "*** RELINK completed in ${THISDIR} ***"

test: functest
	@touch fail
	./functest
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
