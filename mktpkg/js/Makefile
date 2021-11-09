DIST=../../dist/wreis
js:
	mkdir -p res
	cat elems/utils.js elems/image.js elems/table.js elems/jb.js > res/core.js
	echo "completed"
clean:
	./mpak.sh -c
	rm -f res/core.js
	echo "completed"
test:
	echo "completed"
package:
	mkdir -p "${DIST}/mpak/res"
	cp res/core.js res/template00.ai "${DIST}/mpak/res/"
	cp mpak.sh "${DIST}/mpak/"
	echo "completed"

db:
	/usr/local/bin/mysql --no-defaults wreis < mpt.sql

check:
	shellcheck mpak.sh
