TOPDIR="../.."
DIST=wreis
THISDIR=schema

schema:
	@touch fail
	id
	whoami
	cat schema.sql  > db.sql
	mysql --no-defaults < db.sql
	@rm -rf fail
	@echo "*** completed in ${THISDIR} ***"

clean:
	rm -f db.sql fail
	@echo "*** CLEAN completed in ${THISDIR} ***"

test:
	@echo "*** TEST completed in ${THISDIR} ***"

package:
	cp schema.sql ../../dist/${DIST}
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@echo "*** SECURE completed in ${THISDIR} ***"
