#!/bin/bash
#
# USAGE:
# 	getallconf.sh
#
# SYNOPSIS: download all the config files from the repos
#
# DESCRIPTION:
#	As we externalize more things, it is necessary to update all
#	the config files.  This is a convenience script to get them all.
# 
# Usage
#   $1

getfile.sh accord/db/confdev.json
getfile.sh accord/db/confprod.json
getfile.sh accord/db/conflocal.json
getfile.sh accord/db/confwreisprod.json
getfile.sh accord/db/configwreisdev.json
