DIRS=util session db csv ws server webui test
.PHONY: test

wreis:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	rm -rf dist

test:
	for dir in $(DIRS); do make -C $$dir test;done

package:
	for dir in $(DIRS); do make -C $$dir package;done

all: clean wreis package test stats
	echo "Completed"

build: clean wreis package

release:
	/usr/local/accord/bin/release.sh wreis

snapshot:
	cd dist ; rm -f wreis.tar.gz ; tar cvfz wreis.tar.gz wreis ; cd ..

stats:
	@echo
	@echo "-------------------------------------------------------------------------------"
	@echo "GO SOURCE CODE STATISTICS"
	@find . -name "*.go" | srcstats
	@echo "-------------------------------------------------------------------------------"
	@echo "JAVASCRIPT"
	@wc -l dist/wreis/static/js/wreis.js
	@echo "-------------------------------------------------------------------------------"
	@cat test/testreport.txt
