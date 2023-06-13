#!/bin/bash

if [[ $1 == *".."* ]]; then
  echo "El directorio no puede tener .."
  exit 1
fi

HOST="$2"
USER="$3"


ssh "$USER"@"$HOST" "mkdir -p $1"
ssh "$USER"@"$HOST" "mkdir -p $1/users/"
ssh "$USER"@"$HOST" "chmod 770 $1"
ssh "$USER"@"$HOST" "chmod 770 $1/users/"
ssh "$USER"@"$HOST" "chown :uzantoj $1/users/"
ssh "$USER"@"$HOST" "chown :uzantoj $1"


exit 0
