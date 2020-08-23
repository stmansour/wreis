#!/bin/bash
MAJVER=1
MINVER=0
BLDNOFILE="./buildno"
BUILDNO=$(cat ${BLDNOFILE})
BUILDNO=$((BUILDNO + 1))
FNAME="ver.go"
VER=$(printf "%d.%d.%06d" ${MAJVER} ${MINVER} ${BUILDNO})
BLD="${HOSTNAME}"
DAT=$(date)
cat >${FNAME} <<ZZ1EOF
package ws
// THIS FILE IS AUTOGENERATED
// DO NOT EDIT

// GetVersionNo returns the version string
func GetVersionNo() string { return "${VER}" }

// GetBuildMachine returns the name of the machine on which this program was built
func GetBuildMachine() string { return "${BLD}" }

// GetBuildTime returns the timestamp when this program was built
func GetBuildTime() string { return "${DAT}" }
ZZ1EOF
echo "${BUILDNO}" >${BLDNOFILE}