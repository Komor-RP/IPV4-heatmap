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


SELECT addr.latitude, addr.longitude, addr.frequency, log(addr.frequency)/log(2) FROM addresses addr

INNER JOIN 
    (
        SELECT latitude, MAX(frequency) AS sda
        FROM addresses
        GROUP BY latitude
    ) 

WHERE latitude < 100 AND latitude > -100
AND longitude > -100 AND longitude < 0


SELECT latitude, longitude, log(frequency) / log()
FROM addresses WHERE
latitude < ($1) AND latitude > ($2)
AND longitude > ($3) AND longitude < ($4)