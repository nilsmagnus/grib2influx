GRIB2 Golang parser application and library
================================

Parser and library for grib2 file format. 

Forked from github.com/analogic/grib which is now abandoned by the author (see comment on pull request [https://github.com/analogic/grib/pull/1]).

## Usage

Install by typing

    go get -u github.com/nilsmagnus/grib


### Application Usage:

    $ grib -h 
    
    Usage of grib:
     -category int
       	Filters on Category within discipline. -1 means all categories (default -1)
     -dataExport
       	Export data values. (default true)
     -discipline int
       	Filters on Discipline. -1 means all disciplines (default -1)
     -export int
       	Export format. Valid types are 0 (none) 1(print discipline names) 2(print categories) 3(json) 
     -file string
       	Grib filepath
     -latMax int
       	Maximum latitude multiplied with 100000. (default 36000000)
     -latMin int
       	Minimum latitude multiplied with 100000.
     -longMax int
       	Maximum longitude multiplied with 100000. (default 9000000)
     -longMin int
       	Minimum longitude multiplied with 100000. (default -9000000)
     -maxmsg int
       	Maximum number of messages to parse. Does not work in combination with filters. (default 2147483647)
     -operation string
       	Operation. Valid values: 'parse', 'reduce'. (default "parse")
     -reducefile string
       	Destination for reduced file. (default "reduced.grib2")

#### Examples:

Reduce input file to default output-file with discipline 0 (Meteorology):

    grib -operation reduce -file testdata/reduced.grib2 -discipline 0

Filter on area on size of norway+sweden, output to json:
      
    grib -file testdata/gfs.t00z.pgrb2.2p50.f003  -latMin 57000000 -latMax 71000000 -longMin 4400000 -longMax 32000000 -export 3

Filter on temperature only:

    grib -file testdata/gfs.t00z.pgrb2.2p50.f003 -discipline 0 -category 0 

## Library examples

Have a look at 'main.go' for main usage. 

## What works?

- basic binary parsing of GRIB2 GFS files from NOAA
- implemented only "Grid point data - complex packing and spatial differencing"

## TODOs

- Support different types of grids, not only grid0
- Support different types of products, not only product0
- Tests for reduction
- Tests for reading all sections

## Help appreciated

Feel free to fork and submit pull requests or simply create issues for improvements :)

# Grib Documentation

Grib specification:

http://www.wmo.int/pages/prog/www/WMOCodes/Guides/GRIB/GRIB2_062006.pdf

Documentation from noaa.gov :

http://www.nco.ncep.noaa.gov/pmb/docs/on388/


Examples can be found at

http://www.ftp.ncep.noaa.gov/data/nccf/com/gfs/prod/
