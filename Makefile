DIRS=util session db csv ws server admin webui test mktpkg
DIST=dist
.PHONY: test

wreis:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	rm -rf dist

db1:
	cd test/ws; mysql wreis < xh.sql ; cd ../..

test:
	for dir in $(DIRS); do make -C $$dir test;done
	cd test/ws; mysql wreis < xh.sql ; cd ../..

package:
	for dir in $(DIRS); do make -C $$dir package;done

all: clean wreis package test stats
	echo "Completed"

build: clean wreis package

release:
	/usr/local/accord/bin/release.sh wreis

tarzip:
	cd ${DIST};if [ -f ./wreis/config.json ]; then mv ./wreis/config.json .; fi
	cd ${DIST};rm -f wreis.tar*;tar czf wreis.tar.gz wreis
	cd ${DIST};if [ -f ./config.json ]; then mv ./config.json ./wreis/config.json; fi

snapshot: tarzip
	cd ${DIST}; /usr/local/accord/bin/snapshot.sh wreis.tar.gz

stats:
	@echo
	@echo "-------------------------------------------------------------------------------"
	@echo "GO SOURCE CODE STATISTICS"
	@find . -name "*.go" | srcstats
	@echo "-------------------------------------------------------------------------------"
	@echo "JAVASCRIPT"
	@wc -l ${DIST}/wreis/static/js/wreis.js
	@echo "-------------------------------------------------------------------------------"
	@cat test/testreport.txt
