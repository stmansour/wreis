TOP=..
BINDIR=${TOP}/dist/wreis
THISDIR="server"
COUNTOL=${TOP}/tools/bashtools/countol.sh

server: *.go
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	go vet
	@if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

package: server
	mkdir -p ${BINDIR}
	cp server ${BINDIR}/wreis
	cp config.json ${BINDIR}/

build: server package

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f ver.go tmp *.log
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean server package

secure:
	@rm -f config.json configwreisdev.json confprod.json
