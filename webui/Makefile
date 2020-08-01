TOP=..
BINDIR=${TOP}/dist/wreis
STATIC=${BINDIR}/static
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
	mkdir -p ${STATIC}
	cp -r html ${STATIC}/
	cp html/*.css ${STATIC}/html/
	mkdir -p ${STATIC}/js
	cp js/*.js ${STATIC}/js/
	cp js/*.css ${STATIC}/js/
	cp js/elems/*.js ${STATIC}/js
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
