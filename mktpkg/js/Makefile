DIST=../../dist/wreis
js:
	cat elems/utils.js elems/image.js elems/table.js elems/jb.js > core.js
	echo "completed"
clean:
	./mpak.sh -c
	rm -f core.js
	echo "completed"
test:
	echo "completed"
package:
	mkdir -p "${DIST}/mpak"
	cp core.js mpak.sh template00.ai "${DIST}/mpak/"
	echo "completed"

db:
	/usr/local/bin/mysql --no-defaults wreis < mpt.sql

check:
	shellcheck mpak.sh
