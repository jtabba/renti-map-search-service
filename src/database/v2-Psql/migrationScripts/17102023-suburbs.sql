CREATE TABLE Suburbs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    postcode INTEGER NOT NULL,
    geolocation POINT NOT NULL,
    geocode JSON NOT NULL
);