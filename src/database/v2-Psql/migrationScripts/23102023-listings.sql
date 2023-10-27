CREATE TABLE listings(
    id SERIAL PRIMARY KEY,
    address VARCHAR(300) NOT NULL,
    suburb_id INTEGER NOT NULL,
    geolocation GEOGRAPHY(Point) NOT NULL,
    geocode JSON NOT NULL,
    images TEXT [],
    baths INTEGER NOT NULL,
    beds INTEGER NOT NULL,
    parking INTEGER,
    agency JSON,
    weekly_price VARCHAR(100) NOT NULL,
    suburb VARCHAR(100) NOT NULL,
    state VARCHAR(50) NOT NULL,
);

ALTER TABLE listings
    ADD COLUMN postcode VARCHAR(4) NOT NULL,
    ADD COLUMN country VARCHAR(50) NOT NULL,
    ADD COLUMN inspection_open_time VARCHAR(100) NOT NULL,
    ADD COLUMN inspection_close_time VARCHAR(100) NOT NULL,
    ADD COLUMN property_type VARCHAR(50) NOT NULL;

ALTER TABLE listings
    ADD CONSTRAINT fk_suburb_id FOREIGN KEY (suburb_id) REFERENCES aus_suburbs(id),
    -- ADD CONSTRAINT unique_listing UNIQUE (address, suburb_id);