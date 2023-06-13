#!/bin/bash

if [[ $1 == *".."* ]]; then
  echo "El directorio no puede tener .."
  exit 1
fi

HOST="$2"
USER="$3"

ssh "$USER"@"$HOST" "mkdir -p $1"
ssh "$USER"@"$HOST" "chmod 770 $1"
ssh "$USER"@"$HOST" "chmod 770 $1"

# El usuario de sistema es igual al nombre del directorio. Por ejemplo:
# Directorio: mi_dominio_ar
# Usuario: mi_dominio_ar

#ssh uzantoj_create_user@mi-server.local "chown $1:uzantoj $1"
ssh "$USER"@"$HOST" "chown $1:uzantoj $1"

exit 0
