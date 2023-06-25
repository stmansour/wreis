TOP=..
BINDIR=${TOP}/dist/wreis
THISDIR="ws"
COUNTOL=${TOP}/tools/bashtools/countol.sh

ws: ver.go *.go
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	go vet
	golint
	staticcheck
	go build
	go install
	@rm -f fail
	@echo "*** ${THISDIR}: completed ws ***"

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f ver.go tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean ws package

secure:
	@rm -f config.json confdev.json confprod.json

package:
	@echo "*** ${THISDIR}: completed package ***"
