TOP=..
BINDIR=${TOP}/dist/wreis
THISDIR="lib"
COUNTOL=${TOP}/tools/bashtools/countol.sh

lib: *.go
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	go install
	@rm -f fail
	@echo "*** ${THISDIR}: completed lib ***"

clean:
	go clean
	rm -f ver.go tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean lib package

secure:
	@rm -f config.json confdev.json confprod.json

package:
	@echo "*** ${THISDIR}: completed package ***"
