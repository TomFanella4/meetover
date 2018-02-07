#!/bin/bash

APP=/bin/ra-backend
log_file=backend_script.log

curr_time=$(date "+%Y.%m.%d-%H.%M.%S")
echo "[$curr_time] ra-backend Started.." >> $log_file

until $APP; do
	curr_time=$(date "+%Y.%m.%d-%H.%M.%S")
    echo "[$curr_time] ra-backend crashed with exit code $?.  Respawning.." >> $log_file
    # ping the switch?
    sleep 1
done

curr_time=$(date "+%Y.%m.%d-%H.%M.%S")
echo "**********[$curr_time] ra-backend crashed with exit code 0**********" >> $log_file