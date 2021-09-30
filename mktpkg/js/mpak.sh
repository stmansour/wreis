#!/usr/bin/env bash


LOGFILE="log"
REQCOUNT=0
COOKIES=
USER="smansour"
PRID=1
OUTFILE="jbx.js"
PROPJSON="property.json"
CWD=$(pwd)

HOST="http://localhost:8276"
#HOST="https://showponyinvestments.com"

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
	${CMD} | tee serverreply | python -m json.tool >${3} 2>>${LOGFILE}
}

#-----------------------------------------------------------------------------
# Clean - remove old files first...
#-----------------------------------------------------------------------------
Clean() {
    rm -f request response loginrequest log serverreply "${OUTFILE}" property.json portfolio.ai "${PROPJSON}" Img*
}

#-----------------------------------------------------------------------------
# Login...
#-----------------------------------------------------------------------------
LIReq() {

    if [ "${PASSWD}x" = "x" ]; then
        echo "Your password is required to access the WREIS server."
        echo "You can enter it at the prompt, or to avoid having to enter it"
        echo "you can export it in an environment variable as follows:"
        echo "    PASSWD=\"your password\""
        echo "    export PASSWD"
        read -sp 'Password: ' PASSWD
    fi
    encodeRequest "{\"user\":\"${USER}\",\"pass\":\"${PASSWD}\"}"   # puts encoded request in file named "request"
    dojsonPOST "${HOST}/v1/authn/" "request" "response"  # URL, JSONfname, serverresponse

    #-----------------------------------------------------------------------------
    # Now we need to add the token to the curl command for future calls to
    # the server.  curl -b "air=${TOKEN}"  ...
    # Set the command line for cookies in ${COOKIES} and dojsonPOST will use them.
    #-----------------------------------------------------------------------------
    TOKEN=$(grep Token "response" | awk '{print $2;}' | sed 's/[",]//g')
    if [ "${TOKEN}x" == "x" ]; then
        echo "Login failed"
        exit 1
    fi
    COOKIES="-b air=${TOKEN}"   # COOKIES is used by dojsonPOST()
}

#-----------------------------------------------------------------------------
# Read propertu PRID
#-----------------------------------------------------------------------------
ReadProperty () {
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "${HOST}/v1/property/${PRID}" "request" "response"  # URL, JSONfname, serverresponse
    cat response | sed 's/^[{}]$//' | sed 's/^[     ]*"record":/var property = /' | grep -v '"status":' | sed 's/},/};/' > "${PROPJSON}"
}

#-----------------------------------------------------------------------------
# BuildJS
#-----------------------------------------------------------------------------
BuildJS () {
    cat header.js > "${OUTFILE}"
    echo "jb.cwd = \"${CWD}\";" >> "${OUTFILE}"
    cat "${PROPJSON}" utils.js jb.js >> "${OUTFILE}"
}

#-----------------------------------------------------------------------------
# GetImages
#-----------------------------------------------------------------------------
GetImages () {

    for (( i = 1; i < 9; i++ )); do
        iname=$(echo "Img${i}" | sed 's/ *//g')
        iurl=$(grep "${iname}" ${PROPJSON} | awk '{print $2}' | sed 's/[",]//g')
        if [ "${iurl}x" != "x" ]; then
            echo -n "downloading image ${iurl}... "
            fname=$(basename -- "${iurl}")
            ext="${fname##*.}"
            curl -s "${iurl}" -o "${iname}.${ext}"
            echo "  ${iname}.${ext} completed"
        fi
    done
}

###############################################################################
###############################################################################

while getopts "c" o; do
	echo "o = ${o}"
	case "${o}" in
	c)	Clean
		echo "cleaned temporary files"
        exit 0
		;;
    *)  echo "Unrecognized option:  ${o}"
        exit 1
        ;;
    esac
done
shift $((OPTIND-1))

Clean
LIReq
echo -n "Pulling information for Property (PRID): ${PRID}... "
ReadProperty
GetImages
echo "Done"
echo "Generating script to create Adobe Illustrator marketing package... "
BuildJS
echo "Done"

exit 0
