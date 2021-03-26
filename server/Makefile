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
# $ jfrog rt u[pload] configwreisdev.json accord/misc/
# $ deployfile.sh configwreisdev.json accord/db/

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/configwreisdev.json
	cp configwreisdev.json config.json

pushconf: configwreisdev.json
	jfrog rt u configwreisdev.json accord/misc/
	/usr/local/accord/bin/deployfile.sh configwreisdev.json accord/db/

package: server
	mkdir -p ${BINDIR}
	cp server ${BINDIR}/wreis
	cp configwreisdev.json ${BINDIR}/config.json

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

secure:
	@rm -f config.json configwreisdev.json confprod.json
