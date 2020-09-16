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
# ALTER TABLE RenewOption MODIFY Opt VARCHAR(100) NOT NULL DEFAULT '';
# CREATE TABLE Traffic (
#     TID BIGINT NOT NULL AUTO_INCREMENT,                     -- A Traffic ID
#     PRID BIGINT NOT NULL DEFAULT 0,                         -- Associated Property
#     FLAGS BIGINT NOT NULL DEFAULT 0,                        --
#     Count BIGINT NOT NULL DEFAULT 0,                        -- number of vehicles per day, or whatever - see Description
#     Description VARCHAR(128) NOT NULL DEFAULT '',           -- Describes Count
#     LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
#     LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
#     CreateTS TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- when was this record created
#     CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
#     PRIMARY KEY (TID)
# );
# ALTER TABLE Property ADD FlowState BIGINT NOT NULL DEFAULT 0 AFTER RSLID;

#=====================================================
#  Put modifications to schema in the lines below
#=====================================================

cat > "${MODFILE}" << LEOF
ALTER TABLE Property ADD Img1 VARCHAR(2048) NOT NULL DEFAULT '' AFTER HQCountry;
ALTER TABLE Property ADD Img2 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img1;
ALTER TABLE Property ADD Img3 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img2;
ALTER TABLE Property ADD Img4 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img3;
ALTER TABLE Property ADD Img5 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img4;
ALTER TABLE Property ADD Img6 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img5;
LEOF

#=====================================================
#  Put dir/sqlfilename in the list below
#=====================================================
declare -a dbs=(
	'ws/xb.sql'
)

for f in "${dbs[@]}"
do
	echo "DROP DATABASE IF EXISTS wreis; CREATE DATABASE wreis; USE wreis; GRANT ALL PRIVILEGES ON wreis.* TO 'ec2-user'@'localhost';" | ${MYSQL}
	echo -n "${f}: loading... "
	${MYSQL} ${DBNAME} < ${f}
	echo -n "updating... "
	${MYSQL} ${DBNAME} < ${MODFILE}
	echo -n "saving... "
	${MYSQLDUMP} ${DBNAME} > ${f}
	echo "done"
done
