#!/bin/bash
PORT=8276
CHECKINGPERIOD=10                  # seconds
LOGFILE="wreiswatchdog.log"

#----------------------------------------------------
#  Main loop:
#  Ping wreis on localhost:8275.
#  If we don't hear back, then restart
#----------------------------------------------------
while true
do
        R=$(curl -s http://localhost:${PORT}/v1/ping | grep -c "WREIS - Version")
        if [ 0 = "${R}" ]; then
                (echo -n "Ping to wreis failed at "; date; echo -n "Restart...") >> ${LOGFILE}

                PID=$(ps -ef | grep "wreis$" | sed 's/[^ ][^ ]* //' | sed 's/[ ].*//')
                if [ "x${PID}" != "x" ]; then
                        echo "kill -9 ${PID}"
                        kill -9 "${PID}"
                fi
                ./activate.sh start
        fi

    #---------------------------------------------------------------------------
    # Touch the logfile, so we know that this process is active.
    # The timestamp on ${LOGFILE} shows when the process was last
    # checked.
    # Wait for a bit, then do it all again...
    #---------------------------------------------------------------------------
    touch "${LOGFILE}"
    sleep ${CHECKINGPERIOD}
done
