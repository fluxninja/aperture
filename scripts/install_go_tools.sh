#!/usr/bin/env bash

tools=$(grep _ ./tools/tools.go | awk -F'"' '{print $2}')
# $tools contains list of tools separated by new line. Use parallel command to execute "go install {}"for each tool in $tools.
parallel --bar --eta --no-notice "go install {}" ::: "$tools"
