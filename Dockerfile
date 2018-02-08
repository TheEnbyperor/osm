FROM ubuntu:14.04

RUN apt-get update && apt-get install -y curl && curl -sL https://deb.nodesource.com/setup_8.x | bash - && apt-get remove -y curl && apt-get install -y make python g++ libboost-all-dev libfreetype6-dev libxml2-dev libpng12-dev libcairo-dev libtiff4-dev libproj-dev libgdal-dev libcurl4-openssl-dev ttf-unifont curl nodejs && npm -g install carto

WORKDIR /src
COPY mapnik mapnik
RUN cd mapnik && make install && ldconfig && cd ..

COPY osm-bright-master osm-bright-master
RUN cd osm-bright-master && ./make.py && cd ../OSMBright/ && carto project.mml > OSMBright.xml

COPY osm ./
