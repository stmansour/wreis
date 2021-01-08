TOP=..
BINDIR=${TOP}/dist/wreis
THISDIR="admin"
COUNTOL=${TOP}/tools/bashtools/countol.sh

admin:
	@rm -f fail
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

package:
	mkdir -p ${BINDIR}
	cp *.sh ${BINDIR}/

build: admin package

clean:
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean server package

secure:
	@rm -f config.json confdev.json confprod.json
