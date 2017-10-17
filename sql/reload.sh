#!/bin/bash

dropdb makeict
createdb makeict -O $1

psql -U $1 -h localhost makeict -f schema.sql
