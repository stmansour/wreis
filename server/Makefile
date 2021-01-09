TOP=..
BINDIR=${TOP}/dist/wreis
THISDIR="server"
COUNTOL=${TOP}/tools/bashtools/countol.sh

server: *.go config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

# you can update the config files as follows:
# $ jfrog rt u[pload] confdev.json accord/misc/
# $ deployfile.sh confdev.json accord/db/

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

pushconf: confdev.json
	jfrog rt u confdev.json accord/misc/
	/usr/local/accord/bin/deployfile.sh confdev.json accord/db/

package: server
	mkdir -p ${BINDIR}
	cp server ${BINDIR}/wreis
	cp confdev.json ${BINDIR}/config.json

build: server package

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f ver.go tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean server package

mkprod:
	/usr/local/accord/bin/getfile.sh accord/db/confprod.json
	cp confprod.json config.json
secure:
	@rm -f config.json confdev.json confprod.json
