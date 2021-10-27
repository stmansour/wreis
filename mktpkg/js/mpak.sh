#!/usr/bin/env bash

LOGFILE="log"
REQCOUNT=0
COOKIES=
PRID=0
OUTFILE="jbx.js"
PROPJSON="property.json"
ROPTJSON="ropt.json"
RENTJSON="rent.json"
SKIPIMAGES=0
CWD=$(pwd)

HOST="http://localhost:8276"
HOST="https://showponyinvestments.com"

ShowPlan() {

    cat << EOF
*************************************************************************
             Server: ${HOST}
               User: ${USERNAME}
    Property (PRID): ${PRID}
*************************************************************************
EOF
}

Usage() {
    cat <<FEOF
mpak.sh

DESCRIPTION
    mpak.sh is a shell script to create an Adobe Illustrator script that
    produces the WREIS Marketing Package based on the Property ID.

USAGE
    mpak.sh [OPTIONS]

    OPTIONS:

    -c	Clean. Removes any temporary files in the directory.

    -p  PRID
        PRID specifies the Property ID for which the marketing package will
        be created.  It must be a number greater than 0.

    -s  Causes the images to NOT be downloaded. Only use this option if you
        really know what you are doing.

    -u  Display this usage writeup.

Examples
    ./mpak.sh -p 34

FEOF
}

#------------------------------------------------------------------------------
#  encodeRequest is just like encodeURI except that it saves the output
#      into a file named "request"
#
#  INPUTS
#  $1  The string to encode
#
#  RETURNS
#      nothing, but the encoded string will be in a file named "request"
#------------------------------------------------------------------------------
encodeRequest() {
  local string="${1}"
  local strlen=${#string}
  local encoded=""
  local pos c o

  for (( pos=0 ; pos<strlen ; pos++ )); do
     c=${string:$pos:1}
     case "$c" in
        [-_.~a-zA-Z0-9] ) o="${c}" ;;
        * )               printf -v o '%%%02x' "'$c"
     esac
     encoded+="${o}"
  done
  echo "${encoded}" > request
}

########################################################################
# dojsonPOST()
#   Simulate a POST command to the server and use
#   the supplied file name as the json data.
#
#	Parameters:
# 		$1 = url
#       $2 = json file
# 		$3 = base file name
########################################################################
dojsonPOST () {
	((REQCOUNT++))
	COOK=""
	if [ "${COOKIES}x" != "x" ]; then
		COOK="${COOKIES}"
	fi
	CMD="curl ${COOK} -s -X POST ${1} -H \"Content-Type: application/json\" -d @${2}"
	${CMD} | tee serverreply | python -m json.tool > "${3}" 2>>${LOGFILE}
}

#-----------------------------------------------------------------------------
# Clean - remove old files first...
#-----------------------------------------------------------------------------
Clean() {
    rm -f request response loginrequest log serverreply "${OUTFILE}" property.json portfolio.ai "${PROPJSON}" "${ROPTJSON}" "${RENTJSON}"
    if [ ${SKIPIMAGES} -eq 0 ]; then
        rm -f Img*
    fi
}

#-----------------------------------------------------------------------------
# Login...
#-----------------------------------------------------------------------------
LIReq() {
    DONE=0
    if [ "${USERNAME}x" = "x" ]; then
        echo "Your username is required to access the WREIS server."
        echo "You can enter it at the prompt, or to avoid having to enter it"
        echo "you can export it in an environment variable as follows:"
        echo "    USERNAME=\"your username\""
        echo "    export USERNAME"
    else
        DONE=1
    fi
    while (( DONE == 0 )); do
        read -rp 'username: ' USERNAME
        if (( ${#USERNAME} < 1 )); then
            echo "come on now, you gotta give me something..."
        else
            DONE=1
        fi
    done

    if [ "${PASSWD}x" = "x" ]; then
        echo "Your password is required to access the WREIS server."
        echo "You can enter it at the prompt, or to avoid having to enter it"
        echo "you can export it in an environment variable as follows:"
        echo "    PASSWD=\"your password\""
        echo "    export PASSWD"
        read -rsp 'password: ' PASSWD
        echo
        echo "got it."
    fi

    DONE=0
    if (( PRID == 0 )); then
        echo "You must supply a Property ID (PRID) greater than 0"
    else
        DONE=1
    fi
    while [ ${DONE} -eq 0 ]; do
        if [ ${PRID} -eq 0 ]; then
            read -p 'PRID: ' ptmp
            if [[ ${ptmp} =~ ^[0-9]+$ ]]; then
                if (( ptmp < 1)); then
                    echo "the PRID must be greater than 0"
                else
                    PRID=${ptmp}
                    DONE=1
                fi
            else
                echo "you must enter a number"
            fi
        fi
    done

    echo -n "Logging into server... "
    encodeRequest "{\"user\":\"${USERNAME}\",\"pass\":\"${PASSWD}\"}"   # puts encoded request in file named "request"
    dojsonPOST "${HOST}/v1/authn/" "request" "response"  # URL, JSONfname, serverresponse

    #-----------------------------------------------------------------------------
    # Now we need to add the token to the curl command for future calls to
    # the server.  curl -b "air=${TOKEN}"  ...
    # Set the command line for cookies in ${COOKIES} and dojsonPOST will use them.
    #-----------------------------------------------------------------------------
    TOKEN=$(grep Token "response" | awk '{print $2;}' | sed 's/[",]//g')
    if [ "${TOKEN}x" == "x" ]; then
        echo
        echo "Login failed. Check your username and password and try again."
        exit 1
    fi
    COOKIES="-b air=${TOKEN}"   # COOKIES is used by dojsonPOST()
    echo "successfully logged in"
}

#-----------------------------------------------------------------------------
# Read property PRID
#-----------------------------------------------------------------------------
GetProperty () {
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "${HOST}/v1/property/${PRID}" "request" "response"  # URL, JSONfname, serverresponse
    ERR=$(grep "status" < response | grep -c "error")
    if (( ERR == 1 )); then
        echo "*** SERVER REPLIED WITH AN ERROR ***"
        grep "message" <response | sed 's/"//g' | sed 's/  *message: //' | sed 's/\\n,//'
        exit 1
    fi
    sed 's/^[{}]$//' <response | sed 's/^[     ]*"record":/var property = /' | grep -v '"status":' | sed 's/},/};/' > "${PROPJSON}"

    GetImages

    # we need RSLID and ROLID
    RSLID=$(grep "RSLID" property.json | sed 's/^[^:][^:]*: //' | sed 's/,//')
    ROLID=$(grep "ROLID" property.json | sed 's/^[^:][^:]*: //' | sed 's/,//')
}

#-----------------------------------------------------------------------------
# GetImages
#-----------------------------------------------------------------------------
GetImages () {
    if [ ${SKIPIMAGES} -eq 0 ]; then
        for (( i = 1; i < 9; i++ )); do
            iname=$(echo "Img${i}" | sed 's/ *//g')
            iurl=$(grep "${iname}" ${PROPJSON} | sed 's/^  *"I..[0-9][0-9]*...//' | sed 's/[",]//g')
            if [ "${iurl}x" != "x" ]; then
                echo -n "[img${i}]"
                fname=$(basename -- "${iurl}")
                ext="${fname##*.}"
                url=$(echo "${iurl}" | sed 's/ /%20/g')
                curl -s "${url}" -o "${iname}.${ext}"
            fi
        done
    fi
}

#-----------------------------------------------------------------------------
# Read RenewOptions
#-----------------------------------------------------------------------------
GetRenewOptions () {
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "${HOST}/v1/renewoptions/${ROLID}" "request" "response"  # URL, JSONfname, serverresponse
    sed 's/^[{}]$//' < response | sed 's/^[     ]*"record":/var property = /' | grep -v '"status":' | sed 's/    ],/    ];/' | sed "s/\"records\":/property[\"renewOptions\"] = /" > "${ROPTJSON}"
}

#-----------------------------------------------------------------------------
# Read RentSteps
#-----------------------------------------------------------------------------
GetRentSteps () {
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "${HOST}/v1/rentsteps/${RSLID}" "request" "response"  # URL, JSONfname, serverresponse
    sed 's/^[{}]$//' < response | sed 's/^[     ]*"record":/var property = /' | grep -v '"status":' | sed 's/    ],/    ];/' | sed "s/\"records\":/property[\"rentSteps\"] = /" > "${RENTJSON}"
}

#-----------------------------------------------------------------------------
# BuildJS
#-----------------------------------------------------------------------------
BuildJS () {
    cat > "${OUTFILE}" <<FFEOF
//
//  jbx.js - the portfolio writer :-)
//
//  File -> Scripts -> Other Script...
//     or press Cmd + F12
//  then select this file.
//=========================================================================
var jb = {
    portfolio: null,        // the portfolio.ai we are auto-generating
    ab: null,               // active artboard
    doc: null,              // the working document
    chattr: null,           // the default font and attributes
    cwd: "",                // the current working directory
    subjProp: 6,            // index of first subject property after cover photo
    lotSizeLabels: [        // what units for LotSize
        "sqft", "acres"
        ],
    ownershipTypeLabels: [      // OwnershipTypetype
        "Fee Simple",
        "Leasehold"
        ],
    ownershipLabels: [
        "Private",
        "Public"
        ],
    roofStructureLabels: [      // roof responsibility
        "Tenant Responsible",
        "Landlord Responsible"
        ],
    guarantorLabels: [      // who is guarantor
        "Corporate",
        "Franchise",
        "Individual"
        ],
    leaseTypeLabels: [
        "Absolute NNN",
        "Double Net",
        "Triple Net",
        "Gross"
    ],
};
FFEOF
    echo "jb.cwd = \"${CWD}\";" >> "${OUTFILE}"
    cat "${PROPJSON}" "${ROPTJSON}" "${RENTJSON}" core.js >> "${OUTFILE}"
}

###############################################################################
###############################################################################

while getopts "csp:u" o; do
	# echo "o = ${o}"
	case "${o}" in
	c)	Clean
		echo "cleaned temporary files"
        exit 0
		;;
    p)  PRID="${OPTARG}"
        echo "PRID set to ${PRID}"
        ;;
    s)  SKIPIMAGES=1
        echo "do not load images"
        ;;
    u)  Usage
        exit 0
        ;;
    *)  echo "Unrecognized option:  ${o}"
        Usage
        exit 1
        ;;
    esac
done
shift $((OPTIND-1))

Clean       # Remove any old files
LIReq       # Log in
ShowPlan
GetProperty
GetRenewOptions
GetRentSteps
BuildJS
echo
echo "Finished"
echo "Execute Adobe Illustrator script named ${OUTFILE}"

exit 0
