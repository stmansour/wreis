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

#------------------------------------------------------------------------------
#  login - will attempt to login to the wreis server. If it is successful
#          it will set two environment variables:
#
#          TOKEN   - will contain the cookie value for AIR login
#          COOKIES - contains the option for CURL to include the AIR cookie
#                    in requests
#
#          dojsonPOST is setup to use ${COOKIES}
#
#  Scenario:
#  Execute the url to ping the server
#
#  Expected Results:
#   1.  It should return the server version
#------------------------------------------------------------------------------

login() {
    if [ "x${COOKIES}" = "x" ]; then
        encodeRequest "{\"user\":\"${USER}\",\"pass\":\"${PASS}\"}"
        OUTFILE="loginrequest"
        dojsonPOST "http://localhost:8276/v1/authn/" "request" "${OUTFILE}"  "login"
        #-----------------------------------------------------------------------------
        # Now we need to add the token to the curl command for future calls to
        # the server.  curl -b "air=${TOKEN}"  ...
        # Set the command line for cookies in ${COOKIES} and dojsonPOST will use them.
        #-----------------------------------------------------------------------------
        TOKEN=$(grep Token "${OUTFILE}" | awk '{print $2;}' | sed 's/[",]//g')
        COOKIES="-b air=${TOKEN}"
    fi
}

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
#  Validate that commands requiring a session will not operate without
#  a session cookie
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
    mysql --no-defaults wreis < x${TFILES}.sql
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/property/" "request" "${TFILES}${STEP}"  "Property-Search"
fi

#------------------------------------------------------------------------------
#  TEST c
#
#  Login and Get a specific property then logoff.
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

    # encodeRequest "{\"user\":\"${USER}\",\"pass\":\"${PASS}\"}"
    # OUTFILE="${TFILES}${STEP}"
    # dojsonPOST "http://localhost:8276/v1/authn/" "request" "${OUTFILE}"  "Property-Search"
    #
    # #-----------------------------------------------------------------------------
    # # Now we need to add the token to the curl command for future calls to
    # # the server.  curl -b "air=${TOKEN}"  ...
    # # Set the command line for cookies in ${COOKIES} and dojsonPOST will use them.
    # #-----------------------------------------------------------------------------
    # TOKEN=$(grep Token "${OUTFILE}" | awk '{print $2;}' | sed 's/[",]//g')
    # COOKIES="-b air=${TOKEN}"
    login

    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/property/" "request" "${TFILES}${STEP}"  "Property-Search"

    encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8276/v1/logoff/" "request" "${TFILES}${STEP}"  "logoff"
    COOKIES=

    #-----------------------------------------------------------------------------
    # At this point, the cookie is no longer valid.  But try to use it again to
    # verify that the server won't let it be used, and it properly handles a
    # terminated cookie value
    #-----------------------------------------------------------------------------
    encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8276/v1/logoff/" "request" "${TFILES}${STEP}"  "logoff using invalid cookie token"
fi

#------------------------------------------------------------------------------
#  TEST d
#
#  Read the rentsteps for RSLID 4
#  Write rentsteps
#
#  Scenario:
#  login
#  read rentsteps
#
#  Expected Results:
#   1. Expecting 3 rent step items
#   2. Write 4 rent steps back.  Only 1 change (added a new one)
#------------------------------------------------------------------------------
TFILES="d"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < xb.sql
    login
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/rentsteps/4" "request" "${TFILES}${STEP}"  "RentSteps"

    # Add one, remove one, and change one...
    encodeRequest '{"cmd":"save","records":[{"recid":6,"RSID":6,"RSLID":4,"Dt":"1/1/2018","Opt":"Year 1","Rent":3500,"FLAGS":0},{"recid":8,"RSID":8,"RSLID":4,"Dt":"1/1/2020","Opt":"Year 3","Rent":3300,"FLAGS":0},{"recid":-1,"RSLID":0,"RSID":-1,"Opt":"asdfasdf","Dt":"Wed, 01 Jan 2020 08:00:00 GMT","Rent":3333,"FLAGS":0}]}'
    dojsonPOST "http://localhost:8276/v1/rentsteps/4" "request" "${TFILES}${STEP}"  "RentSteps"

    # Change all of them to DATE based, and change 1 date
    encodeRequest '{"cmd":"save","records":[{"recid":6,"RSID":6,"RSLID":4,"Dt":"1/1/2018","Opt":"Year 1","Rent":3500,"FLAGS":1},{"recid":8,"RSID":8,"RSLID":4,"Dt":"1/1/2020","Opt":"Year 3","Rent":3300,"FLAGS":1},{"recid":-1,"RSLID":0,"RSID":-1,"Opt":"asdfasdf","Dt":"Wed, 15 Jan 2020 08:00:00 GMT","Rent":3333,"FLAGS":1}]}'
    dojsonPOST "http://localhost:8276/v1/rentsteps/4" "request" "${TFILES}${STEP}"  "RentSteps"
fi

#------------------------------------------------------------------------------
#  TEST e
#
#  Read the renewoptions for ROLID 4
#  Write renewoptions
#
#  Scenario:
#  login
#  read renewoptions
#
#  Expected Results:
#   1. Expecting 3 rent step items
#   2. Write 4 rent steps back.  Only 1 change (added a new one)
#------------------------------------------------------------------------------
TFILES="e"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < xb.sql
    login
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/renewoptions/1" "request" "${TFILES}${STEP}"  "RenewOption"

    # Add one, remove one, and change one...
    encodeRequest '{"cmd":"save","records":[{"recid":1,"ROID":1,"ROLID":1,"Dt":"7/4/2024","Opt":"1","Rent":109709.45,"FLAGS":1},{"recid":3,"ROID":3,"ROLID":1,"Dt":"11/4/2026","Opt":"3","Rent":114141.71,"FLAGS":1},{"recid":-1,"ROID":-1,"ROLID":1,"Dt":"12/15/2027","Opt":"3","Rent":20000.00,"FLAGS":1}]}'
    dojsonPOST "http://localhost:8276/v1/renewoptions/1" "request" "${TFILES}${STEP}"  "RenewOptions"

    # Change all of them to DATE based, and change 1 date
    encodeRequest '{"cmd":"save","records":[{"recid":1,"ROID":1,"ROLID":1,"Dt":"7/4/2024","Opt":"Year 1","Rent":109709.45,"FLAGS":0},{"recid":3,"ROID":3,"ROLID":1,"Dt":"11/4/2026","Opt":"Year 2","Rent":114141.71,"FLAGS":0},{"recid":-1,"ROID":-1,"ROLID":1,"Dt":"12/15/2027","Opt":"Year 3","Rent":20000.00,"FLAGS":0}]}'
    dojsonPOST "http://localhost:8276/v1/renewoptions/1" "request" "${TFILES}${STEP}"  "RenewOptions"
fi

stopWsrv
echo "WREIS SERVER STOPPED"

logcheck

exit 0
