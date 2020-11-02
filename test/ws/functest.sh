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

#------------------------------------------------------------------------------
#  TEST f
#
#  Save property form
#
#  Scenario:
#  login
#  save an updated version of the entry where we remove the apostrophe from
#  the name
#
#  Expected Results:
#   1. Expecting 3 rent step items
#   2. Write 4 rent steps back.  Only 1 change (added a new one)
#------------------------------------------------------------------------------
TFILES="f"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < xb.sql
    login

    # write it
    encodeRequest '{"cmd":"save","record":{"recid":4,"FlowState":1,"PID":0,"PRID":4,"Name":"Sallys Sludge Salon","YearsInBusiness":5,"ParentCompany":"","URL":"https://bbe.com/","Symbol":"","Price":77777.88,"DownPayment":510000,"RentableArea":16000,"RentableAreaUnits":0,"LotSize":26,"LotSizeUnits":1,"CapRate":0.28,"AvgCap":0.35,"BuildDate":"Wed, 01 Jan 1975 08:00:00 GMT","FLAGS":2,"Ownership":1,"TenantTradeName":"Sallys Sludge Salon","LeaseGuarantor":1,"LeaseType":2,"DeliveryDt":"Wed, 01 Jan 1975 08:00:00 GMT","OriginalLeaseTerm":31,"RentCommencementDt":"Fri, 15 Jun 2018 07:00:00 GMT","LeaseExpirationDt":"Mon, 15 Jun 2020 07:00:00 GMT","TermRemainingOnLease":71,"TermRemainingOnLeaseUnits":0,"ROLID":2,"RSLID":3,"Address":"1235 Elm Street","Address2":"","City":"Goober","State":"AK","PostalCode":"12345","Country":"USA","LLResponsibilities":"","NOI":25000,"HQAddress":"1235 Elm Street","HQAddress2":"","HQCity":"Goober","HQState":"AK","HQPostalCode":"12345","HQCountry":"USA","CreateTime":"1900-01-01 00:00:00 UTC","CreateBy":0,"LastModTime":"1900-01-01 00:00:00 UTC","LastModBy":0}}'
    dojsonPOST "http://localhost:8276/v1/property/4" "request" "${TFILES}${STEP}"  "Property-save"

    # read it back to make sure the changes stuck
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/property/4" "request" "${TFILES}${STEP}"  "Read_property"

fi

#------------------------------------------------------------------------------
#  TEST g
#
#  Save Traffic
#
#  Scenario:
#  login
#  Read the traffic info from the db
#
#  Expected Results:
#   1. Expecting 3 rent step items
#   2. Write 4 rent steps back.  Only 1 change (added a new one)
#------------------------------------------------------------------------------
TFILES="g"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < xb.sql
    login

    # read it back to make sure the changes stuck
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/trafficitems/1" "request" "${TFILES}${STEP}"  "Read_traffic"

    # write it - change the count on two
    encodeRequest '{"cmd":"save","records": [{"Count": 9999,"Description": "Vehicles per day on Main street","FLAGS": 0,"PRID": 1,"TID": 1,"recid": 1},{"Count": 7777,"Description": "Elm Street","FLAGS": 0,"PRID": 1,"TID": 2,"recid": 2}]}'
    dojsonPOST "http://localhost:8276/v1/trafficitems/1" "request" "${TFILES}${STEP}"  "traffic-save"

    # read it back, make sure the change stuck
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/trafficitems/1" "request" "${TFILES}${STEP}"  "Read_traffic"

    # remove one, add a new one
    encodeRequest '{"cmd":"save","records": [{"Count": 9999,"Description": "Vehicles per day on Main street","FLAGS": 0,"PRID": 1,"TID": 1},{"Count": 13458,"Description": "Parade Ave","FLAGS": 0,"PRID": 1,"TID": -1}]}'
    dojsonPOST "http://localhost:8276/v1/trafficitems/1" "request" "${TFILES}${STEP}"  "traffic-save"

    # read it back, verify the changes
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/trafficitems/1" "request" "${TFILES}${STEP}"  "Read_traffic"

    # add a new one
    encodeRequest '{"cmd":"save","records":[{"recid":1,"TID":1,"PRID":1,"Description":"Vehicles per day on Main street","Count":9999,"FLAGS":0},{"recid":3,"TID":3,"PRID":1,"Description":"Parade Ave","Count":13458,"FLAGS":0},{"recid":-2,"PRID":1,"TID":-2,"Description":"Odd Street","Count":899,"FLAGS":0}]}'
    dojsonPOST "http://localhost:8276/v1/trafficitems/1" "request" "${TFILES}${STEP}"  "traffic-save"

    # read it back, verify the changes
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/trafficitems/1" "request" "${TFILES}${STEP}"  "Read_traffic"
fi
#------------------------------------------------------------------------------
#  TEST h
#
#  State Info
#
#  Scenario:
#
#  Expected Results:
#------------------------------------------------------------------------------
TFILES="h"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < xh.sql
    login
    # read what we have
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/1" "request" "${TFILES}${STEP}"  "Read-StateInfo"

    # save a new one, modify one.  New one added is SIID 2
    encodeRequest '{"cmd":"save","records":[{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":1,"InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"PRID":1,"SIID":1,"recid":1},{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":2,"InitiatorDt":"2020-10-01 10:23:45 UTC","InitiatorUID":211,"PRID":1,"SIID":-1,"recid":2}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/1" "request" "${TFILES}${STEP}"  "Save-StateInfo"

    # add a third
    encodeRequest '{"cmd":"save","records":[{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":1,"InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"PRID":1,"SIID":1,"recid":1},{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":2,"InitiatorDt":"2020-10-01 10:23:45 UTC","InitiatorUID":211,"PRID":1,"SIID":2,"recid":2},{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":3,"InitiatorDt":"2020-10-02 10:23:45 UTC","InitiatorUID":211,"PRID":1,"SIID":-2,"recid":2}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/1" "request" "${TFILES}${STEP}"  "Save-StateInfo"

    # read to make sure we have 3
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/1" "request" "${TFILES}${STEP}"  "Read-StateInfo"

    # remove the third
    encodeRequest '{"cmd":"save","records":[{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":1,"InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"PRID":1,"SIID":1,"recid":1},{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":2,"InitiatorDt":"2020-10-01 10:23:45 UTC","InitiatorUID":211,"PRID":1,"SIID":2,"recid":2}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/1" "request" "${TFILES}${STEP}"  "Save-StateInfo"

    # read to make sure we have 2
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/1" "request" "${TFILES}${STEP}"  "Read-StateInfo"

    # read property 7 to make sure all info is being correctly fetched
    encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/7" "request" "${TFILES}${STEP}"  "Read-StateInfo"
fi
#------------------------------------------------------------------------------
#  TEST i
#
#  State Info:
#       process a Reject request
#
#  Scenario:
#
#  Expected Results:
#------------------------------------------------------------------------------
TFILES="i"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    mysql --no-defaults wreis < xh.sql
    login

    # 0. Error case 1
    # try save a reject that we're not listed as the Authorizer
    encodeRequest '{"cmd":"reject","records":[{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":3,"Reason":"This is the reason","InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"PRID":3,"SIID":6,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "StateInfo-Reject-error-case-1"

    # 1. Error case 2
    # try save a reject where we changed the UID to our current UID, but it does not match the one in the database
    encodeRequest '{"cmd":"reject","records":[{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":269,"FLAGS":0,"FlowState":3,"Reason":"This is the reason","InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"PRID":3,"SIID":6,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "StateInfo-Reject-error-case-2"

    # 2. Error case 3
    # try to save a reject without a reason
    encodeRequest '{"cmd":"reject","records":[{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":269,"FLAGS":0,"FlowState":4,"InitiatorDt":"2020-10-04 10:37:45 UTC","InitiatorUID":92,"PRID":4,"SIID":10,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/4" "request" "${TFILES}${STEP}"  "StateInfo-Reject-error-case-3"

    # 3. save a reject
    encodeRequest '{"cmd":"reject","records":[{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":269,"FLAGS":0,"FlowState":4,"Reason":"Listing pictures look bad","InitiatorDt":"2020-10-04 10:37:45 UTC","InitiatorUID":92,"PRID":4,"SIID":10,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/4" "request" "${TFILES}${STEP}"  "StateInfo-Reject"

    # 4. read property 4 to make sure all info is being correctly fetched
    encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/4" "request" "${TFILES}${STEP}"  "Read-StateInfo"

    # 5. set ready status - Error case - only owner can set READY
    encodeRequest '{"cmd":"ready","records":[{"PRID": 4,"Reason": "","SIID": 33,"ApproverDt": "1900-01-01 00:00:00 UTC","ApproverName": "William Tester","ApproverUID": 269,"CreateBy": 269,"CreateByName": "William Tester","CreateTime": "2020-10-30 22:29:08 UTC","FLAGS": 0,"FlowState": 4,"InitiatorDt": "2020-10-31 00:00:00 UTC","InitiatorName": "Patrick Long","InitiatorUID": 92,"LastModBy": 269,"LastModByName": "William Tester","LastModTime": "2020-10-30 22:29:08 UTC","recid": 33}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/4" "request" "${TFILES}${STEP}"  "StateInfo-Ready-Error1"

    # 6. set ready status - this should work as tester is the owner of the task
    encodeRequest '{"cmd":"ready","records":[{"PRID": 3,"Reason": "","SIID": 6,"ApproverDt": "1900-01-01 00:00:00 UTC","ApproverUID": 80,"CreateBy": 269,"CreateByName": "William Tester","CreateTime": "2020-10-30 22:29:08 UTC","FLAGS": 0,"FlowState": 4,"InitiatorDt": "2020-10-31 00:00:00 UTC","InitiatorName": "Patrick Long","InitiatorUID": 92,"LastModBy": 269,"LastModByName": "William Tester","LastModTime": "2020-10-30 22:29:08 UTC","recid": 33}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "StateInfo-Ready"

    # 7. read property 3 to make sure all info is being correctly fetched.  We just set the state to READY.  FLAGS should be
    encodeRequest '{"cmd":"get"}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "Read-StateInfo"

    # 8. save an approval.  This should create SIID 34, with FLAGS=0 and FlowState = 5
    encodeRequest '{"cmd":"approve","records":[{"PRID": 4,"Reason": "","SIID": 33,"ApproverDt": "1900-01-01 00:00:00 UTC","ApproverName": "William Tester","ApproverUID": 269,"CreateBy": 269,"CreateByName": "William Tester","CreateTime": "2020-10-30 22:29:08 UTC","FLAGS": 0,"FlowState": 4,"InitiatorDt": "2020-10-31 00:00:00 UTC","InitiatorName": "Patrick Long","InitiatorUID": 92,"LastModBy": 269,"LastModByName": "William Tester","LastModTime": "2020-10-30 22:29:08 UTC","recid": 33}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/4" "request" "${TFILES}${STEP}"  "StateInfo-Approve"

    # 9. save an approval - should cause an error - the work was already completed
    encodeRequest '{"cmd":"approve","records":[{"PRID": 4,"Reason": "","SIID": 33,"ApproverDt": "1900-01-01 00:00:00 UTC","ApproverName": "William Tester","ApproverUID": 269,"CreateBy": 269,"CreateByName": "William Tester","CreateTime": "2020-10-30 22:29:08 UTC","FLAGS": 0,"FlowState": 4,"InitiatorDt": "2020-10-31 00:00:00 UTC","InitiatorName": "Patrick Long","InitiatorUID": 92,"LastModBy": 269,"LastModByName": "William Tester","LastModTime": "2020-10-30 22:29:08 UTC","recid": 33}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/4" "request" "${TFILES}${STEP}"  "StateInfo-Approve-Error"

    # 10. save revert but don't include a reason
    encodeRequest '{"cmd":"revert","records":[{"PRID": 4,"Reason": "","SIID": 34,"ApproverDt": "1900-01-01 00:00:00 UTC","ApproverName": "William Tester","ApproverUID": 269,"CreateBy": 269,"CreateByName": "William Tester","CreateTime": "2020-10-30 22:29:08 UTC","FLAGS": 0,"FlowState": 4,"InitiatorDt": "2020-10-31 00:00:00 UTC","InitiatorName": "Patrick Long","InitiatorUID": 92,"LastModBy": 269,"LastModByName": "William Tester","LastModTime": "2020-10-30 22:29:08 UTC","recid": 33}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/4" "request" "${TFILES}${STEP}"  "StateInfo-Revert-Error-case-1"

    # 11. save revert
    encodeRequest '{"cmd":"revert","records":[{"PRID": 4,"Reason": "Picture shows the wrong elevation of the building","SIID": 34,"ApproverDt": "1900-01-01 00:00:00 UTC","ApproverName": "William Tester","ApproverUID": 269,"CreateBy": 269,"CreateByName": "William Tester","CreateTime": "2020-10-30 22:29:08 UTC","FLAGS": 0,"FlowState": 4,"InitiatorDt": "2020-10-31 00:00:00 UTC","InitiatorName": "Patrick Long","InitiatorUID": 92,"LastModBy": 269,"LastModByName": "William Tester","LastModTime": "2020-10-30 22:29:08 UTC","recid": 33}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/4" "request" "${TFILES}${STEP}"  "StateInfo-Revert"

    # 12. try save a revert that we're not listed as the Authorizer
    encodeRequest '{"cmd":"revert","records":[{"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":3,"Reason":"This is the reason","InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"PRID":3,"SIID":6,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "StateInfo-Revert-error-case-2"

    # 13. try to revert something that's in state 1
    encodeRequest '{"cmd":"revert","records":[{"PRID":1,"SIID":1,"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":269,"FLAGS":0,"FlowState":1,"Reason":"This is the reason","InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/1" "request" "${TFILES}${STEP}"  "StateInfo-Revert-error-case-3"

    # 14. change the owner on a property where we're not the owner or approver
    encodeRequest '{"cmd":"setowner","records":[{"OwnerUID":47,"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":203,"FLAGS":0,"FlowState":3,"Reason":"We need to get this moving","InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"PRID":3,"SIID":6,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "StateInfo-setowner"

    # 15. change the approver on a property where we're not the owner or approver - this is an error case as it tries to use a finished state info
    encodeRequest '{"cmd":"setapprover","records":[{"PRID":3,"SIID":6,"OwnerUID":47,"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":72,"FLAGS":0,"FlowState":3,"Reason":"Someone needs to approve this","InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "StateInfo-setapprover-error-case-1"

    # 16. change the approver on a property where we're not the owner or approver - this should work since we are operating on the latest state info
    encodeRequest '{"cmd":"setapprover","records":[{"PRID":3,"SIID":36,"OwnerUID":47,"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":72,"FLAGS":0,"FlowState":3,"Reason":"Someone needs to approve this","InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "StateInfo-setapprover"

    # 17. this should fail as we completed work on SIID 36 in the last command
    encodeRequest '{"cmd":"setapprover","records":[{"PRID":3,"SIID":36,"OwnerUID":47,"ApproverDt":"1970-01-01 00:00:00 UTC","ApproverUID":72,"FLAGS":0,"FlowState":3,"Reason":"Someone needs to approve this","InitiatorDt":"2020-10-01 10:37:45 UTC","InitiatorUID":211,"recid":1}]}'
    dojsonPOST "http://localhost:8276/v1/stateinfo/3" "request" "${TFILES}${STEP}"  "StateInfo-setapprover-error"

fi

stopWsrv
echo "WREIS SERVER STOPPED"

logcheck

exit 0
