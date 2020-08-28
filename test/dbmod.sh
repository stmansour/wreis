#!/bin/bash

#==========================================================================
#  This script performs SQL schema changes on the test databases that are
#  saved as SQL files in the test directory. It loads them, performs the
#  ALTER commands, then saves the sql file.
#
#  If the test file uses its own database saved as a .sql file, make sure
#  it is listed in the dbs array
#==========================================================================

MODFILE="dbqqqmods.sql"
MYSQL="mysql --no-defaults"
MYSQLDUMP="mysqldump --no-defaults"
DBNAME="wreis"

#=====================================================
#  Retain prior changes as comments below
#=====================================================
# ALTER TABLE Property MODIFY LeaseCommencementDt RentCommencementDt DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00';
# ALTER TABLE Property ADD TermRemainingOnLeaseUnits SMALLINT NOT NULL DEFAULT 0 AFTER TermRemainingOnLease;
# ALTER TABLE RenewOption DROP COLUMN Count;

#=====================================================
#  Put modifications to schema in the lines below
#=====================================================

cat > "${MODFILE}" << LEOF
ALTER TABLE RenewOption MODIFY Opt VARCHAR(100) NOT NULL DEFAULT '';
LEOF

#=====================================================
#  Put dir/sqlfilename in the list below
#=====================================================
declare -a dbs=(
	'ws/xb.sql'
)

for f in "${dbs[@]}"
do
	echo -n "${f}: loading... "
	${MYSQL} ${DBNAME} < ${f}
	echo -n "updating... "
	${MYSQL} ${DBNAME} < ${MODFILE}
	echo -n "saving... "
	${MYSQLDUMP} ${DBNAME} > ${f}
	echo "done"
done
