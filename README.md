Work in progress

# grib2influx

cli tool to parse forecasts (or other sections) from grib-files to influx as timeseries

# Detailed description


* parse the grib2-file using the griblib, a golang library to read grib2-files
* create points from the data and insert them into influx. 
   * series-name is deducted from the coordinate of the datapoint
   * values are added with the respective category name

From one single grib-file you will end up with one time-point for each series with as many values as categories in the grib-file(wind, temperature etc)


# TODO 
figure out forecast-time versus value-time description, "wide table" design


