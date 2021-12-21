DIST=../../dist/wreis
LDIST=mpak
js:
	mkdir -p res
	cat elems/utils.js elems/image.js elems/table.js elems/jb.js > res/core.js
	echo "completed"
clean:
	./mpak.sh -c
	rm -rf res/core.js "${LDIST}"
	echo "completed"
test:
	echo "completed"
package:
	mkdir -p "${LDIST}/res"
	cp res/template00.ai "${LDIST}/res/"
	cp res/core.js  "${LDIST}/res/"
	cp mpak.sh "${LDIST}/"
	tar cvf mpak.tar mpak ; gzip mpak.tar
	cp mpak.tar.gz "${DIST}/static/"
	rm -rf mpak.tar.gz mpak
	echo "completed"

db:
	/usr/local/bin/mysql --no-defaults wreis < mpt.sql

check:
	shellcheck mpak.sh
