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
# ALTER TABLE Property ADD Img1 VARCHAR(2048) NOT NULL DEFAULT '' AFTER HQCountry;
# ALTER TABLE Property ADD Img2 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img1;
# ALTER TABLE Property ADD Img3 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img2;
# ALTER TABLE Property ADD Img4 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img3;
# ALTER TABLE Property ADD Img5 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img4;
# ALTER TABLE Property ADD Img6 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img5;
# ALTER TABLE Property ADD Img7 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img6;
# ALTER TABLE Property ADD Img8 VARCHAR(2048) NOT NULL DEFAULT '' AFTER Img7;

# ALTER TABLE Property CHANGE CreateTS CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
# ALTER TABLE RenewOption CHANGE CreateTS CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
# ALTER TABLE RenewOptions CHANGE CreateTS CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
# ALTER TABLE RentStep CHANGE CreateTS CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
# ALTER TABLE RentSteps CHANGE CreateTS CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

# CREATE TABLE StateInfo (
#     SIID BIGINT NOT NULL AUTO_INCREMENT,                    -- State Info ID
#     PRID BIGINT NOT NULL DEFAULT 0,                         -- Associated Property
#     FLAGS BIGINT NOT NULL DEFAULT 0,                        --
#     FlowState BIGINT NOT NULL DEFAULT 0,                    --
#     InitiatorUID BIGINT NOT NULL DEFAULT 0,                 --
#     ApproverUID BIGINT NOT NULL DEFAULT 0,                 --
#     LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
#     LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
#     CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,-- when was this record created
#     CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
#     PRIMARY KEY (SIID)
# )
# ALTER TABLE StateInfo ADD InitiatorDt DATE NOT NULL DEFAULT '1970-01-01 00:00:00' AFTER InitiatorUID;
# ALTER TABLE StateInfo ADD Reason VARCHAR(256) NOT NULL DEFAULT '' AFTER ApproverDt;

# ALTER TABLE StateInfo CHANGE InitiatorUID OwnerUID BIGINT NOT NULL DEFAULT 0;
# ALTER TABLE StateInfo CHANGE InitiatorDt OwnerDt DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00';

# ALTER TABLE Property CHANGE LotSize LotSize DECIMAL(19,4) NOT NULL DEFAULT 0;

# ALTER TABLE Property CHANGE YearsInBusiness YearBuilt SMALLINT NOT NULL DEFAULT 0;
# ALTER TABLE Property CHANGE YearBuilt YearFounded SMALLINT NOT NULL DEFAULT 0;
# ALTER TABLE Property DROP COLUMN HQAddress, DROP COLUMN HQAddress2, DROP COLUMN HQPostalCode, DROP COLUMN HQCountry;
# ALTER TABLE Property DROP COLUMN DeliveryDt;

# ALTER TABLE Property CHANGE BuildDate BuildYear SMALLINT NOT NULL DEFAULT 0;
# ALTER TABLE Property ADD RenovationYear SMALLINT NOT NULL DEFAULT 0 AFTER BuildYear;

# ALTER TABLE Property DROP COLUMN TermRemainingOnLease, DROP COLUMN TermRemainingOnLeaseUnits;

# ALTER TABLE Property CHANGE Ownership OwnershipType SMALLINT NOT NULL DEFAULT 0;

#=====================================================
#  Put modifications to schema in the lines below
#=====================================================

cat > "${MODFILE}" << LEOF
ALTER TABLE Property CHANGE Img1 Img1 VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE Property CHANGE Img2 Img2 VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE Property CHANGE Img3 Img3 VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE Property CHANGE Img4 Img4 VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE Property CHANGE Img5 Img5 VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE Property CHANGE Img6 Img6 VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE Property CHANGE Img7 Img7 VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE Property CHANGE Img8 Img8 VARCHAR(256) NOT NULL DEFAULT '';

ALTER TABLE Property ADD Img9 VARCHAR(256) NOT NULL DEFAULT '' AFTER Img8;
ALTER TABLE Property ADD Img10 VARCHAR(256) NOT NULL DEFAULT '' AFTER Img9;
ALTER TABLE Property ADD Img11 VARCHAR(256) NOT NULL DEFAULT '' AFTER Img10;
ALTER TABLE Property ADD Img12 VARCHAR(256) NOT NULL DEFAULT '' AFTER Img11;
ALTER TABLE RenewOptions ADD MPText VARCHAR(256) NOT NULL DEFAULT '' AFTER ROLID;
ALTER TABLE RentSteps ADD MPText VARCHAR(256) NOT NULL DEFAULT '' AFTER RSLID;
LEOF

#=====================================================
#  Put dir/sqlfilename in the list below
#=====================================================
declare -a dbs=(
	'ws/xb.sql'
	'ws/xh.sql'
	'photo/xa.sql'
	'../mktpkg/js/mpt.sql'
	'../mktpkg/samples/mktdb.sql'
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
