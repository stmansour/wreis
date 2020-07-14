TOP=..
BINDIR=${TOP}/dist/wreis
THISDIR="ws"
COUNTOL=${TOP}/tools/bashtools/countol.sh

ws: ver.go *.go config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed mojosrv ***"

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f ver.go tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean ws package

mkprod:
	/usr/local/accord/bin/getfile.sh accord/db/confprod.json
	cp confprod.json config.json
secure:
	@rm -f config.json confdev.json confprod.json
