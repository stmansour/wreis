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

function CreateNewDB() {
	${MYSQL} -h "${HOST}" -P "${PORT}" <schema.sql
}

function MakeProdDB() {
        MYSQLOPTS=

        HOST=$(grep "WREISDbhost" config.json | sed -e 's/.*"WREISDbhost"[ \t]*:[ \t]*"//' | sed -e 's/",$//')
        PORT=$(grep "WREISDbport" config.json | sed -e 's/.*"WREISDbport"[ \t]*:[ \t]*//' | sed -e 's/[ \t]*,[ \t]*$//')

        # echo "HOST = ${HOST}, PORT = ${PORT}"

        MYSQLDUMP="${MYSQL}dump"
		if [ -f "${DBBACKUP}" ]; then
			echo "Production database exists, and backup file (${DBBACKUP}) also exists."
			echo "Please rename backup file or remove it and run this script again if you want to empty the database."
			exit 1
		fi
        echo "execute: ${MYSQLDUMP} -h ${HOST} -P ${PORT} ${DBNAME} >${DBBACKUP}"
        ${MYSQLDUMP} -h "${HOST}" -P "${PORT}" "${DBNAME}" >${DBBACKUP}
        RC=$?
        if [ $RC == 0 ]; then
                echo "WREIS database existed. A backup was made to ${DBBACKUP}."
        else
                # echo "execute: ${MYSQL} -h ${HOST} -P ${PORT} <schema.sql"
                echo "WREIS database was created."
        fi
		CreateNewDB
}

SOURCE="${BASH_SOURCE[0]}"
while [ -h "${SOURCE}" ]; do # resolve ${SOURCE} until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "${SOURCE}" )" && pwd )"
  SOURCE="$(readlink "${SOURCE}")"
  [[ "${SOURCE}" != /* ]] && SOURCE="$DIR/${SOURCE}" # if ${SOURCE} was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "${SOURCE}" )" && pwd )"

pushd "${DIR}" || exit 2
PROD=$(grep '"Env"' config.json | grep -c 1)   # if this is production then PROD == 1, otherwise PROD == 0
if [ "${PROD}" = "1" ]; then
        MakeProdDB
else
        CreateNewDB
fi
popd || exit 2
