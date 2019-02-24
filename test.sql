CREATE TABLE Addresses (
    id SERIAL NOT NULL,
    latitude FLOAT  NOT NULL,
    longitude FLOAT  NOT NULL
);

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
FROM './test.csv' CSV HEADER;

INSERT INTO addresses (latitude, longitude)
SELECT latitude, longitude FROM tmp
WHERE latitude is NOT NULL AND longitude is NOT NULL;