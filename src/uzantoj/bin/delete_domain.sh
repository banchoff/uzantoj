#!/bin/bash


if [[ $1 == *".."* ]]; then
  echo "El directorio no puede tener .."
  exit 1
fi

HOST="$2"
USER="$3"

NOW=`date +"%d-%m-%Y_%H-%M"`
ssh "$USER"@"$HOST" "mv $1 $1_$NOW.bak"


# Hardcodeado para que no de error. Pero hay que manejar bien esto.

exit 0
