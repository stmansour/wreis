js:
	echo "completed"
clean:
	./mpak.sh -c
	echo "completed"
test:
	echo "completed"
package:
	cat elems/utils.js elems/image.js elems/table.js elems/jb.js > core.js
	echo "completed"

db:
	/usr/local/bin/mysql --no-defaults wreis < mpt.sql

check:
	shellcheck mpak.sh
