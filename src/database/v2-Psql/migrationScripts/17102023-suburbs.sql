CREATE TABLE Suburbs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    postcode INTEGER NOT NULL,
    geolocation GEOGRAPHY(Point) NOT NULL,
    geocode JSON NOT NULL,
    CONSTRAINT unique_suburb UNIQUE (name, state, postcode)
);