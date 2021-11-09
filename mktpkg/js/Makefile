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
	mkdir -p "${LDIST}/mpak/res"
	cp res/template00.ai "${LDIST}/mpak/res/"
	cp res/core.js  "${LDIST}/mpak/res/"
	cp mpak.sh "${LDIST}/mpak/"
	tar cvf mpak.tar mpak ; gzip mpak.tar
	cp mpak.tar.gz "${DIST}/"
	rm -rf mpak.tar.gz mpak
	echo "completed"

db:
	/usr/local/bin/mysql --no-defaults wreis < mpt.sql

check:
	shellcheck mpak.sh
