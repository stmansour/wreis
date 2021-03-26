#!/bin/bash
#
# USAGE:
# 	putallconf.sh
#
# SYNOPSIS: download all the config files from the repos
#
# DESCRIPTION:
#	As we externalize more things, it is necessary to update all
#	the config files.  This is a convenience script to write them back
#   to both Responsibility.
#
# Usage
#   $1

deployfile.sh confdev.json accord/db/
jfrog rt u confdev.json accord/misc/
deployfile.sh confprod.json accord/db/
jfrog rt u confprod.json accord/misc/
deployfile.sh conflocal.json accord/db/
jfrog rt u conflocal.json accord/misc/
deployfile.sh confwreisprod.json accord/db/
jfrog rt u confwreisprod.json accord/misc/
deployfile.sh configwreisdev.json accord/db/
jfrog rt u configwreisdev.json accord/misc/
