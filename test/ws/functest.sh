#!/bin/bash

TESTNAME="WebServices"
TESTSUMMARY="Test Wsrv Web Services"
DBGEN=../../../tools/dbgen
CREATENEWDB=0
WBIN="../../dist/wreis"

echo "Create new database..."
mysql --no-defaults wreis < ${WBIN}/schema.sql

source ../../share/base.sh

echo "STARTING WREIS SERVER"
DB2LOADED=0

USER=$(grep Tester1Name config.json | awk '{print $2;}' | sed 's/[",]//g')
PASS=$(grep Tester1Pass config.json | awk '{print $2;}' | sed 's/[",]//g')
if [ "${USER}x" = "x" -o "${PASS}x" = "x" ]; then
    echo "Could not establish user and password. Is config.conf correct?"
    exit 2
fi

startWsrv

#------------------------------------------------------------------------------
#  TEST a
#  ping the server
#
#  Scenario:
#  Execute the url to ping the server
#
#  Expected Results:
#   1.  It should return the server version
#------------------------------------------------------------------------------
TFILES="a"
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    echo "Test ${TFILES}"
    C=$(curl -s http://localhost:8276/v1/ping | tee serverreply | grep "WREIS - Version 1.0" | wc -l | sed 's/ *//' )
    if [ "${C}" = "1" ]; then
        passmsg
    else
        failmsg
    fi
    ((TESTCOUNT++))
fi

#------------------------------------------------------------------------------
#  TEST b
#
#  Validate the property search command
#
#  Scenario:
#  Search
#
#  Expected Results:
#   1.
#   2.
#------------------------------------------------------------------------------
TFILES="b"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    # stopWsrv
    mysql --no-defaults wreis < x${TFILES}.sql
    # startWsrv

    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/property/" "request" "${TFILES}${STEP}"  "Property-Search"
fi

#------------------------------------------------------------------------------
#  TEST c
#
#  Loging and Get a specific property
#
#  Scenario:
#  login
#  get property 1
#
#  Expected Results:
#   1. All fields of property 1 are returned
#   2.
#------------------------------------------------------------------------------
TFILES="c"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < xb.sql

    encodeRequest "{\"user\":\"${USER}\",\"pass\":\"${PASS}\"}"
    OUTFILE="${TFILES}${STEP}"
    dojsonPOST "http://localhost:8276/v1/authn/" "request" "${OUTFILE}"  "Property-Search"

    #-----------------------------------------------------------------------------
    # Now we need to add the token to the curl command for future calls to
    # the server.  curl -b "air=${TOKEN}"  ...
    # Set the command line for cookies in ${COOKIES} and dojsonPOST will use them.
    #-----------------------------------------------------------------------------
    TOKEN=$(grep Token "${OUTFILE}" | awk '{print $2;}' | sed 's/[",]//g')
    echo "Token = ${TOKEN}"
    COOKIES="-b air=${TOKEN}"

    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/property/" "request" "${TFILES}${STEP}"  "Property-Search"
fi

stopWsrv
echo "WREIS SERVER STOPPED"

logcheck

exit 0
