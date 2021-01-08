#!/bin/bash
MYSQLOPTS="--no-defaults"
DBBACKUP="WreisBackupDB.sql"
DBNAME="wreis"
if [ -f /usr/local/bin/mysql ]; then
	MYSQL="/usr/local/bin/mysql"
elif [ -f /usr/bin/mysql ]; then
	MYSQL="/usr/bin/mysql"
else
	MYSQL="mysql"
fi

function MakeProdDB() {
	MYSQLOPTS=

	HOST=$(grep "WREISDbhost" config.json | sed -e 's/.*"WREISDbhost"[ \t]*:[ \t]*"//' | sed -e 's/",$//')
	PORT=$(grep "WREISDbport" config.json | sed -e 's/.*"WREISDbport"[ \t]*:[ \t]*//' | sed -e 's/[ \t]*,[ \t]*$//')
	MYSQLDUMP="${MYSQL}dump"
	${MYSQLDUMP} -h ${HOST} -P ${PORT} ${DBNAME} >${DBBACKUP}
	RC=$?
	if [ $RC == 0 ]; then
		echo "WREIS database existed. A backup was made to ${DBBACKUP}."
	else
		${MYSQL} -h ${HOST} -P ${PORT} <schema.sql
		echo "WREIS database was created."
	fi
}

SOURCE="${BASH_SOURCE[0]}"
while [ -h "${SOURCE}" ]; do # resolve ${SOURCE} until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "${SOURCE}" )" && pwd )"
  SOURCE="$(readlink "${SOURCE}")"
  [[ "${SOURCE}" != /* ]] && SOURCE="$DIR/${SOURCE}" # if ${SOURCE} was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "${SOURCE}" )" && pwd )"

pushd ${DIR}
PROD=$(grep '"Env"' config.json | grep 1 | wc -l)   # if this is production then PROD == 1, otherwise PROD == 0
if [ "${PROD}" = "1" ]; then
	MakeProdDB
else
	mysql ${MYSQLOPTS} <schema.sql
fi
popd