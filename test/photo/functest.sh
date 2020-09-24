#!/bin/bash

TESTNAME="WebServicesPhoto"
TESTSUMMARY="Test Wsrv Photo Mgmt Web Services"
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

        #-----------------------------------------------------------------------
        # This is needed so that the tests can be entered at any point.
        # login() uses dojsonPOST which updates STEP.  We only want the
        # test steps in the main routine below to update the test counts.
        # login should be written so that it can be called anywhere, anytime
        # and it will not alter the sequencing of the output files.
        #-----------------------------------------------------------------------
        ((STEP--))
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
#------------------------------------------------------------------------------
#  TEST a
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
TFILES="a"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < x${TFILES}.sql
    login

    CMD="curl ${COOKIES} -F request={\"cmd\":\"save\",\"PRID\":1,\"idx\":1,\"filename\":\"roller-32.png\"} -F file=@roller-32.png http://localhost:8276/v1/propertyphoto/1/1"
    echo "CMD = ${CMD}"
    ${CMD}

fi

stopWsrv
echo "WREIS SERVER STOPPED"

logcheck

exit 0
