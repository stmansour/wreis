TOP=../..
THISDIR="utiltest"

utiltest: *.go
	@touch fail
	go vet
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	@rm -f fail
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f conf*.json log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? ${THISDIR} fail
	@echo "*** CLEAN completed in ${THISDIR} ***"

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

relink:
	go build
	@echo "*** RELINK completed in ${THISDIR} ***"

test: utiltest
	@touch fail
	./utiltest
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
