CREATE TABLE Addresses (
    id SERIAL NOT NULL,
    latitude FLOAT  NOT NULL,
    longitude FLOAT  NOT NULL
    
);

-- Create temporary table for all the data
CREATE TEMP TABLE tmp (
    extra_column1 TEXT,
    extra_column2 FLOAT,
    extra_column3 FLOAT,
    extra_column4 FLOAT,
    extra_column5 FLOAT,
    extra_column6 FLOAT,
    extra_column7 TEXT,
    latitude FLOAT,
    longitude FLOAT,
    extra_column8 FLOAT
);

\copy tmp
FROM './Geolite2-City-Blocks-IPv4.csv' CSV HEADER;

-- copy over the latitude and longitude columns
-- where latitude and longitude are not null
INSERT INTO addresses (latitude, longitude)
SELECT latitude, longitude FROM tmp
WHERE latitude is NOT NULL AND longitude is NOT NULL;