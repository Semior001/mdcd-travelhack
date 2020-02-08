#!/bin/sh

set -e

echo "executing $@"

host="$1"
shift
cmd="$@"

echo "Trying to authorize with $DB_USER : $DB_PASSWORD to $host ..."

until PGPASSWORD=$DBPASSWORD psql -h "$host" -U "$DBUSER" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd
