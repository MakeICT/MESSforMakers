
REM run this command from the folder containing schema.sql, btree.sql and test_data.sql
REM the command must be run with with super user, port, and owner as arguments e.g. ./reload.cmd postgres_test 5433 db_owner

dropdb -U %1 -p %2 --if-exists makeict
createdb -U %1 -p %2 -O %3 makeict 

psql -U %1 -h localhost -p %2  -f btree.sql makeict
psql -U %3 -h localhost -p %2  -f schema.sql makeict
psql -U %3 -h localhost -p %2  -f test_data.sql makeict
