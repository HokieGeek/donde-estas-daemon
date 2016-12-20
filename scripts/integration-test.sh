#!/bin/bash

here=$(cd ${0%/*}; pwd)
${here}/db.sh

${here}/build.sh &
sleep 40s

smoke.sh
update-person.sh
get-person.sh
