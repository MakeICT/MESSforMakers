#!/bin/bash

dropdb makeict
createdb makeict

psql -U $USER -h localhost makeict -f schema.sql
