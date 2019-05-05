#!/bin/bash

# run this command from the folder containing schema.sql and test_data.sql
# the command must be run with with the postgres user as an argument e.g. ./reload.sh postgres_test

dropdb makeict
createdb makeict -O $1

psql -U $1 -h localhost makeict -f schema.sql
psql -U $1 -h localhost makeict -f test_data.sql
