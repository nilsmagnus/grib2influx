#!/bin/bash
#
ls testdata/gfs* > ._file
while read p; do
    ./grib2influx -gribfile $p

done <._file
rm ._file
echo done