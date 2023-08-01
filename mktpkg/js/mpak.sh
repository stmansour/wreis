#!/usr/bin/env bash

LOGFILE="log"
LOGALL="logall"
CRLFMT="crlfmt.txt"
PERFMON=0
REQCOUNT=0
COOKIES=
PRID=0
OUTFILE="jb.jsx"
PROPJSON="property.json"
ROPTJSON="ropt.json"
RENTJSON="rent.json"
SKIPIMAGES=0
CWD=$(pwd)
SAVECOREFILES=0
FORMATTER="python3" # possible values: "ruby"  "perl"  "python"  "python3"
CHECKTOOLS=0        # by default, we do not check tools

WRHOME="${HOME}/.wreis"
WRCONFIG="${WRHOME}/config"

#HOST="http://localhost:8276"
HOST="https://showponyinvestments.com"

#------------------------------------------------------------------------------
#  CheckTools - A quick check to make sure that tools are in place
#------------------------------------------------------------------------------
CheckTools() {

    local p="0"
    local p3="0"
    notfound=()

    tools=("ruby" "perl" "python" "python3")

    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &>/dev/null; then
            notfound+=("${tool}")
            if [ "${tool}" == "python" ]; then
                p="1"
            elif [ "${tool}" == "python3" ]; then
                p3="1"
            fi
        fi
    done

    #--------------------------------
    # are we missing the formatter?
    #--------------------------------
    for tool in "${notfound[@]}"; do
        if [ "${FORMATTER}" == "${tool}" ]; then
            echo "*** FATAL ERROR ***  the internal formatter (${FORMATTER}) is not on this computer"
            echo "exiting now"
            exit 1
        fi
    done

    #-------------------------------------------------------------
    # return now if we're not going to report missing tools...
    #-------------------------------------------------------------
    if [ "${CHECKTOOLS}" == "0" ]; then
        return
    fi

    #--------------------------------
    #  Report missing tools
    #--------------------------------
    if ((${#notfound[@]} > 0)); then
        echo "The following tools are missing from this system:"
        for tool in "${notfound[@]}"; do
            echo "  ${tool}"
        done
    fi
}

#------------------------------------------------------------------------------
#  ShowPlan - print relevant connection information to the terminal
#------------------------------------------------------------------------------
ShowPlan() {
    cat <<EOF
*************************************************************************
             Server: ${HOST}
               User: ${WUNAME}
    Property (PRID): ${PRID}
          Formatter: ${FORMATTER}
*************************************************************************
EOF
}

#------------------------------------------------------------------------------
#  Usage - call to print instructions to the terminal
#------------------------------------------------------------------------------
Usage() {
    cat <<FEOF
mpak.sh

DESCRIPTION
    mpak.sh is a shell script to create an Adobe Illustrator script that
    produces the WREIS Marketing Package based on the Property ID.  It does
    this by logging into the WREIS server and retrieving the information
    it needs to build the marketing package for a specific property. Then it
    creates a script for Adobe Illustrator in a file named ${OUTFILE}. Open
    Illustrator, then select File -> Scripts -> Other Script... , then select
    ${OUTFILE} . This will create a new tab called portfolio.ai and it will create
    the marketing package based on the data it downloaded.

    In order to log into the server, you will need to provide your username and
    password. The script will ask for these values if it needs them. If you
    want to have these values saved on your system, you can provide the values
    when asked and then indicate that you want them saved.  This information
    will be stored in a special directory in your home directory:
    
        ${HOME}/.wreis
    
    Do this only if you trust the system you are on.

    Alternatively, you can create environment variables for WUNAME and PASSWD
    that the script can use. Here is an example:

        bash$ WUNAME="jsmith"
        bash$ PASSWD="mysecretpassword"
        bash$ export WUNAME PASSWD

USAGE
    mpak.sh [OPTIONS]

    OPTIONS:
    -a  Save temporary files (server response files)

    -c	Clean. Removes any temporary files in the directory.

    -f { ruby | perl | python | python3 }
        Which formatter to use.  Different systems have different tools
        installed, this script can use any of these 4 choices. The
        default is python3.

    -p  PRID
        PRID specifies the Property ID for which the marketing package will
        be created.  It must be a number greater than 0.
        
    -r  Enables curl performance monitoring.

    -s  Causes the images to NOT be downloaded. Only use this option if you
        really know what you are doing.

    -t  Check this system for the existence of the tools needed to
        run this script.

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

    for ((pos = 0; pos < strlen; pos++)); do
        c=${string:$pos:1}
        case "$c" in
        [-_.~a-zA-Z0-9]) o="${c}" ;;
        *) printf -v o '%%%02x' "'$c" ;;
        esac
        encoded+="${o}"
    done
    echo "${encoded}" >request
}

#------------------------------------------------------------------------------
# doformat
#
# INPUTS
# $1 = "1" means use head -1 serverreply,  "2" means cat serverreply | ...
# $2 = output file
#------------------------------------------------------------------------------
doformat() {
    local f

    if [[ "$1" == "1" ]]; then
        f=$(head -n1 serverreply)
    else
        f=$(cat serverreply)
    fi

    case "${FORMATTER}" in
    "ruby")
        echo "${f}" | ruby -rjson -e 'puts JSON.pretty_generate(JSON.parse(ARGF.read))' >"${2}" 2>>${LOGFILE}
        ;;
    "perl")
        echo "${f}" | perl -MJSON -0777 -ne 'print JSON->new->pretty->encode(decode_json($_))' >"${2}" 2>>${LOGFILE}
        ;;
    "python")
        echo "${f}" | python -m json.tool >"${2}" 2>>${LOGFILE}
        ;;
    "python3")
        echo "${f}" | python3 -m json.tool >"${2}" 2>>${LOGFILE}
        ;;
    esac
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
dojsonPOST() {
    ((REQCOUNT++))
    PRF=""
    COOK=""
    if [ "${COOKIES}x" != "x" ]; then
        COOK="${COOKIES}"
    fi
    endpoint=${1}  # first script argument
    json_file=${2} # second script argument

    if [ ${PERFMON} -eq 0 ]; then
        CMD="curl ${COOK} --keepalive-time 2 -s -X POST ${1} -H \"Content-Type: application/json\" -d @${2}"
        ${CMD} >serverreply
        # cat serverreply | python3 -m json.tool >"${3}" 2>>${LOGFILE}
        doformat "2" "${3}"
    else
        CMD="curl ${COOK} -w @${CRLFMT} --keepalive-time 2 -s -X POST ${1} -H \"Content-Type: application/json\" -d @${2}"
        echo "${CMD}" >>"${LOGALL}"
        ${CMD} >serverreply
        doformat "1" "${3}"

        #---------------------
        # python solution
        #---------------------
        # head -1 serverreply | python3 -m json.tool >"${3}" 2>>${LOGFILE}

        #---------------------
        # Perl solution
        #---------------------
        #head -1 serverreply | perl -MJSON -0777 -ne 'print JSON->new->pretty->encode(decode_json($_))' >"${3}" 2>>${LOGFILE}

        #---------------------
        # Ruby solution
        #---------------------
        # head -1 serverreply | ruby -rjson -e 'puts JSON.pretty_generate(JSON.parse(ARGF.read))' > "${3}" 2>>${LOGFILE}

        cat serverreply >>"${LOGALL}"
    fi
}

#-----------------------------------------------------------------------------
# Clean - remove old files first...
#-----------------------------------------------------------------------------
Clean() {
    rm -f request response loginrequest loginresponse log serverreply "${OUTFILE}" property.json portfolio.ai "${PROPJSON}" "${ROPTJSON}" "${RENTJSON}"
    if [ ${SKIPIMAGES} -eq 0 ]; then
        rm -f Img*
    fi
}

#-----------------------------------------------------------------------------
# GetCreds - get the users credentials for login
#-----------------------------------------------------------------------------
GetCreds() {
    if [ "${WUNAME}x" = "x" ]; then
        echo "Your username is required to access the WREIS server."
        # echo "You can enter it at the prompt, or to avoid having to enter it"
        # echo "you can export it in an environment variable as follows:"
        # echo "    WUNAME=\"your username\""
        # echo "    export WUNAME"
    else
        DONE=1
    fi
    while ((DONE == 0)); do
        read -rp 'username: ' WUNAME
        if ((${#WUNAME} < 1)); then
            echo "come on now, you gotta give me something..."
        else
            DONE=1
        fi
    done

    if [ "${PASSWD}x" = "x" ]; then
        echo "Your password is required to access the WREIS server."
        # echo "You can enter it at the prompt, or to avoid having to enter it"
        # echo "you can export it in an environment variable as follows:"
        # echo "    PASSWD=\"your password\""
        # echo "    export PASSWD"
        read -rsp 'password: ' PASSWD
        echo
        echo "got it."
    fi

}

#-----------------------------------------------------------------------------
# SaveLoginInfo - saves login info to ${WRCONFIG}
#-----------------------------------------------------------------------------
SaveLoginInfo() {
    mkdir -p "${WRHOME}"
    if [ ! -f "${WRCONFIG}" ]; then
        cat >"${WRCONFIG}" <<FEOF3
username: ${WUNAME}
password: ${PASSWD}
FEOF3
        chmod 600 "${WRCONFIG}"
    fi
}

#-----------------------------------------------------------------------------
# Login...
#-----------------------------------------------------------------------------
LIReq() {

    #------------------------------------------
    # Read the user login info if we have it
    #------------------------------------------
    if [ -f "${WRCONFIG}" ]; then
        while IFS= read -r line; do
            if [[ $line == username:* ]]; then
                WUNAME="${line#username: }"
            elif [[ $line == password:* ]]; then
                export PASSWD="${line#password: }"
            fi
        done <"${WRCONFIG}"
    else
        GetCreds
    fi

    DONE=0
    if ((PRID == 0)); then
        echo "You must supply a Property ID (PRID) greater than 0"
    else
        DONE=1
    fi

    while [ ${DONE} -eq 0 ]; do
        if [ "${PRID}" -eq 0 ]; then
            read -r -p 'PRID: ' ptmp
            if [[ ${ptmp} =~ ^[0-9]+$ ]]; then
                if ((ptmp < 1)); then
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
    encodeRequest "{\"user\":\"${WUNAME}\",\"pass\":\"${PASSWD}\"}" # puts encoded request in file named "request"
    dojsonPOST "${HOST}/v1/authn/" "request" "response"             # URL, JSONfname, serverresponse

    cat response >loginresponse

    #-----------------------------------------------------------------------------
    # Now we need to add the token to the curl command for future calls to
    # the server.  curl -b "air=${TOKEN}"  ...
    # Set the command line for cookies in ${COOKIES} and dojsonPOST will use them.
    #-----------------------------------------------------------------------------
    # TOKEN=$(grep Token "response" | awk '{print $2;}' | sed 's/[",]//g')  # worked for Python formatter
    TOKEN=$(grep Token "response" | sed -n 's/.*"Token"[ ]*:[ ]*"\(.*\)",/\1/p') # works with Perl formatter

    if [ "${TOKEN}x" == "x" ]; then
        echo
        echo "Login failed. Check your username and password and try again."
        exit 1
    fi

    #-------------------------------------
    # Offer to save login info...
    #-------------------------------------
    if [ ! -f "${WRCONFIG}" ]; then
        while true; do
            read -r -p "Save login information? " answer
            answer=$(echo "$answer" | tr '[:upper:]' '[:lower:]')
            case $answer in
            "yes" | "y")
                SaveLoginInfo
                echo "Saved"
                break
                ;;
            "no" | "n") break ;;
            *) echo "Please enter yes or no." ;;
            esac
        done
    fi

    COOKIES="-b air=${TOKEN}" # COOKIES is used by dojsonPOST()
    echo "successfully logged in"
}

#-----------------------------------------------------------------------------
# Read property PRID
#-----------------------------------------------------------------------------
GetProperty() {
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "${HOST}/v1/property/${PRID}" "request" "response" # URL, JSONfname, serverresponse
    ERR=$(grep "status" <response | grep -c "error")
    if ((ERR == 1)); then
        echo "*** SERVER REPLIED WITH AN ERROR ***"
        grep "message" <response | sed 's/"//g' | sed 's/  *message: //' | sed 's/\\n,//'
        exit 1
    fi
    sed 's/^[{}]$//' <response | sed 's/^[     ]*"record":/var property = /' | grep -v '"status":' | sed 's/},/};/' >"${PROPJSON}"

    GetImages

    # we need RSLID and ROLID
    RSLID=$(grep "RSLID" property.json | sed 's/^[^:][^:]*: //' | sed 's/,//')
    ROLID=$(grep "ROLID" property.json | sed 's/^[^:][^:]*: //' | sed 's/,//')
}

#-----------------------------------------------------------------------------
# GetImages
#-----------------------------------------------------------------------------
GetImages() {
    if [ ${SKIPIMAGES} -eq 0 ]; then
        for ((i = 1; i < 9; i++)); do
            iname=$(echo "Img${i}" | sed 's/ *//g')
            iurl=$(grep "${iname}" ${PROPJSON} | sed 's/^  *"I..[0-9][0-9]*...//' | sed 's/[",]//g')
            if [ "${iurl}x" != "x" ]; then
                echo -n "[img${i}]"
                fname=$(basename -- "${iurl}")
                ext="${fname##*.}"
                url="${iurl// /%20}" # replace spaces with %20
                curl --keepalive-time 2 -s "${url}" -o "${iname}.${ext}"
            fi
        done
    fi
}

#-----------------------------------------------------------------------------
# Read RenewOptions
#-----------------------------------------------------------------------------
GetRenewOptions() {
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "${HOST}/v1/renewoptions/${ROLID}" "request" "response" # URL, JSONfname, serverresponse
    sed 's/^[{}]$//' <response | sed 's/^[     ]*"record":/var property = /' | grep -v '"status":' | sed 's/    ],/    ];/' | sed "s/\"records\":/property[\"renewOptions\"] = /" >"${ROPTJSON}"
}

#-----------------------------------------------------------------------------
# Read RentSteps
#-----------------------------------------------------------------------------
GetRentSteps() {
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "${HOST}/v1/rentsteps/${RSLID}" "request" "response" # URL, JSONfname, serverresponse
    sed 's/^[{}]$//' <response | sed 's/^[     ]*"record":/var property = /' | grep -v '"status":' | sed 's/    ],/    ];/' | sed "s/\"records\":/property[\"rentSteps\"] = /" >"${RENTJSON}"
}

#-----------------------------------------------------------------------------
# BuildJS
#-----------------------------------------------------------------------------
BuildJS() {
    cat >"${OUTFILE}" <<FFEOF
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
    cwd: "${CWD}",          // the current working directory
    subjProp: 6,            // index of first subject property after cover photo
    lotSizeLabels: [        // what units for LotSize
        "SF", "Acres"
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
    ]
};
FFEOF
    cat "${PROPJSON}" "${ROPTJSON}" "${RENTJSON}" res/core.js >>"${OUTFILE}"

    if ((SAVECOREFILES != 1)); then
        echo "removing temporary files..."
        rm -rf "${PROPJSON}" "${ROPTJSON}" "${RENTJSON}" log request response serverreply
    fi

}

SetupPerfmon() {
    cat >"${CRLFMT}" <<FEOF
\n
    time_namelookup:  %{time_namelookup}\n
       time_connect:  %{time_connect}\n
    time_appconnect:  %{time_appconnect}\n
   time_pretransfer:  %{time_pretransfer}\n
      time_redirect:  %{time_redirect}\n
 time_starttransfer:  %{time_starttransfer}\n
                    ----------\n
         time_total:  %{time_total}\n
\n
FEOF
    echo "Created ${CRLFMT}"
    echo "curl Performance Log" >"${LOGALL}"
    date >>"${LOGALL}"
}

###############################################################################
###############################################################################

while getopts "acmtrsf:p:u" o; do
    # echo "o = ${o}"
    case "${o}" in
    a)
        SAVECOREFILES=1
        echo "SAVECOREFILES = ${SAVECOREFILES}"
        ;;
    c)
        Clean
        echo "cleaned temporary files"
        exit 0
        ;;
    f)
        FORMATTER="${OPTARG}"
        case "${FORMATTER}" in
        "ruby" | "perl" | "python" | "python3")
            echo "Formatter is set to ${FORMATTER}"
            ;;
        *)
            echo "Invalid formatter! Please set FORMATTER to either 'ruby', 'perl', 'python', or 'python3'."
            exit 1
            ;;
        esac
        ;;

    r)
        PERFMON=1
        echo "curl performance monitoring enabled. Filename: ${LOGALL}"
        SetupPerfmon
        ;;
    p)
        PRID="${OPTARG}"
        echo "PRID set to ${PRID}"
        ;;
    s)
        SKIPIMAGES=1
        echo "do not load images"
        ;;
    t)
        CHECKTOOLS="1"
        echo "CHECKTOOLS = ${CHECKTOOLS}"
        ;;
    u)
        Usage
        exit 0
        ;;
    *)
        echo "Unrecognized option:  ${o}"
        Usage
        exit 1
        ;;
    esac
done
shift $((OPTIND - 1))

if [ "${1}x" != "x" ]; then
    Usage
    exit 1
fi

CheckTools
Clean # Remove any old files
LIReq # Log in
ShowPlan
GetProperty
GetRenewOptions
GetRentSteps
BuildJS
echo
echo "Finished"
echo "Execute Adobe Illustrator script named ${OUTFILE}"

exit 0
