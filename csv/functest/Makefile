TOP=../..
THISDIR="csv/functest"

functest: *.go config.json
	@touch fail
	@go vet
	golint
	staticcheck
	go build
	@rm -f fail
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f conf*.json log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? functest fail
	@echo "*** CLEAN completed in ${THISDIR} ***"

config.json:
	cp ${TOP}/server/config.json config.json

relink:
	go build
	@echo "*** RELINK completed in ${THISDIR} ***"

test: functest
	@touch fail
	mysql --no-defaults < ../../db/schema/db.sql
	./functest
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
