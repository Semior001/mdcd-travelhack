#!/bin/sh

set -e

echo "executing $@"

host="$1"
shift
cmd="$@"

echo "Trying to authorize with $DB_USER : $DB_PASSWORD to $host ..."

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -U "$DB_USER" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd