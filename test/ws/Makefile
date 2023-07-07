THISDIR="ws"
TOP=../..

all:
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f loginrequest slog wreis.log serverreply log llog *.g ./gold/*.g err.txt [a-z] [a-z][0-9] [a-z0-9][0-9][0-9] qq? fail request w.[a-z] con*.json curlperf
	@echo "*** CLEAN completed in ${THISDIR} ***"

test: config.json
	@touch fail
	./functest.sh
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

config.json:
	cp ${TOP}/server/config.json config.json

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
