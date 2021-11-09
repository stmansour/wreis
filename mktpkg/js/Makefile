DIST=../../dist/wreis
js:
	mkdir -p res
	cat elems/utils.js elems/image.js elems/table.js elems/jb.js > res/core.js
	echo "completed"
clean:
	./mpak.sh -c
	rm -rf res/core.js tmp
	echo "completed"
test:
	echo "completed"
package:
	mkdir -p tmp/mpak/res
	cp res/template00.ai tmp/mpak/res/
	cp res/core.js  tmp/mpak/res/
	cp mpak.sh tmp/mpak/
	cd tmp ; tar cvf mpak.tar . ; gzip mpak.tar
	cp tmp/mpak.tar.gz "${DIST}/"
	echo "completed"

db:
	/usr/local/bin/mysql --no-defaults wreis < mpt.sql

check:
	shellcheck mpak.sh
