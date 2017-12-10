FROM alpine

ADD grib2influx grib2influx

ENTRYPOINT ./grib2influx

EXPOSE 8080