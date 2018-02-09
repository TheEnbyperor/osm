FROM debian:jessie

RUN apt-get update &&\
    apt-get install -y libboost-all-dev git-core build-essential autoconf libtool libgeos-dev libgeos++-dev libpq-dev &&\
    libbz2-dev libproj-dev libgdal-dev liblua5.2-dev lua5.1 liblua5.1-dev &&\
    git clone git://github.com/openstreetmap/osm2pgsql.git &&\
    cd osm2pgsql &&\
    ./autogen.sh &&\
    ./configure &&\
    make &&\
    make install &&\
    rm -rf /var/lib/apt/lists/*

COPY main /
CMD /main