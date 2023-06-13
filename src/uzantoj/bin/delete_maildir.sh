#!/bin/bash


if [[ $1 == *".."* ]]; then
  echo "El directorio no puede tener .."
  exit 1
fi

HOST="$2"
USER="$3"

NOW=`date +"%d-%m-%Y_%H-%M"`
logger "mv $1 $1_$NOW.bak"
ssh "$USER"@"$HOST" "mv $1 $1_$NOW.bak"

exit 0
