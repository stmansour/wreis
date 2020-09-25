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


TFILES="a"
#------------------------------------------------------------------------------
#  TEST a
#
#  Read the rentsteps for RSLID 4
#  Write rentsteps
#
#  Scenario:
#  login
#  store a photo in S3
#  delete the photo from S3
#
#  Expected Results:
#   1. see the steps below
#------------------------------------------------------------------------------
TFILES="a"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < x${TFILES}.sql
    login

    # Send a MimeMultipart for saving and image...
    if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
        CMD="curl -s ${COOKIES} -F request={\"cmd\":\"save\",\"PRID\":1,\"idx\":1,\"filename\":\"roller-32.png\"} -F file=@roller-32.png http://localhost:8276/v1/propertyphoto/1/1"
        echo "CMD = ${CMD}"
        ${CMD} | tee serverreply | python -m json.tool > "${TFILES}${STEP}" 2>>${LOGFILE}
        doCheckOnly "${TFILES}${STEP}"

        # make sure it's there...
        url=$(grep url a0 | sed 's/[^h]*//' | sed 's/"//g')
        echo "url = ${url}"
        img="${TFILES}${STEP}.png"
        curl -s ${url} > ${img}
        df=$(diff a0.png roller-32.png | wc -l | sed 's/ *//')
        if [ ${df} != "0" ]; then
            echo "*** ERROR ***   expecting df = 0, found df = ${df}"
            exit 1
        fi
        rm -rf ${img}
        passmsg

        ((STEP++))
    fi

    if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
        # remove the photo we just stored
        encodeRequest '{"cmd":"delete","PRID":1,"idx":1}'
        dojsonPOST "http://localhost:8276/v1/propertyphotodelete/1/1" "request" "${TFILES}${STEP}"  "DeletePhoto"

        if [ "${url}x" != "x" ]; then
            curl -s ${url} > s3del.txt
            resp=$(cat s3del.txt | grep "<Code>AccessDenied" | wc -l | sed 's/ *//')
            if [ "${resp}" != "1" ]; then
                echo "*** ERROR *** expected response error from s3 in s3del.txt, but no error found!!"
                echo "              cat s3del.txt and see what's going on"
                exit 1
            fi
            rm s3del.txt
            passmsg
        fi
    fi

fi

stopWsrv
echo "WREIS SERVER STOPPED"

logcheck

exit 0
