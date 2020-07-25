TOP=..
BINDIR=${TOP}/dist/wreis
THISDIR="webui"
COUNTOL=${TOP}/tools/bashtools/countol.sh

webui:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	# insert steps here...
	@rm -f fail
	@echo "*** ${THISDIR}: completed webui ***"

jshint:
	@touch fail
	@${COUNTOL} "jshint --extract=always ./html/*.html"
	@rm -rf fail

clean:
	@echo "*** ${THISDIR}: completed clean ***"

package:
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	mkdir -p ${BINDIR}
	cp -r html ${BINDIR}/
	if [ -d "./js" ]; then mkdir -p ${BINDIR}/js;cp js/*.js ${BINDIR}/js/;cp js/elems/*.js ${BINDIR}/js; fi
	@rm -f fail
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean webui package test

images:
	tar cvf images.tar html/images
	gzip wreisimages.tar

secure:
	@rm -f config.json confdev.json confprod.json
