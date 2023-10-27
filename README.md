## Overview

I tasked myself with making an application which uses real-time property data to highlight inspection times of property listings.

This service was set up to retrieve and send data as quickly as possible, and hence Go was chosen to create it (which I had to learn from scratch).

I used the opportunity to learn more about SQL, working with geometry data and Docker containerisation. The pimrary technologies used are:

• PSQL with a Postgis database (no ORM was utilised for increased performance and to force me to manually write queries)
• Docker
• Go with Gin as the HTTP framework of choice due to its performance benefits over the native Go HTTP library
• Go Colly for web scraping
• Air for hot reloading

### Functionality

The service works by - Checking the database for suburbs in the search radius - Checking the cache for existing listings in the database of suburbs within said radius - Scraping + saving what doesn't exist - Sending the data to the front-end

## Running the server

The application uses `Golang 1.21` but is functional with `v1.19 or greater`

1. Download [Go](https://go.dev/doc/install) v1.19 or greater

[air](https://github.com/cosmtrek/air) is used for hot reloading the server

2. Download `air`
3. `air init`
4. `air` to run the server

By default, the server runs on port `8080`

## Running the database

You must have Docker installed to run the application. OrbStack is a useful alternative to Docker UI to help visualise and control existing images and containers

1. Install [Orb](https://orbstack.dev/)

### Start container

Ensure you `cd` into the root of `map-search-service` before proceeding

2. `docker-compose -f docker-compose.yml up`

### Restore database

You will need to acquire a copy of the database schema to restore the necessary tables

3. Get database dump from Jad
4. Run `cat renti-dump.sql | docker exec -i map-search-service-database psql -U postgres -d renti_db` to restore the data

### Connect to container db

Once restored, you can connect to the database container and run a query to confirm the completion of the setup

5. Run `docker compose --env-file /dev/null exec renti-db psql -U postgres -d renti_db`
6. Once connected to `renti_db#=` run

```
    SELECT count(id) from aus_suburbs;
```

Which should output something like

```
     count
    -------
     3657
    (1 row)
```

Congratulations! You are all set up.

### Creating a dump file

If a dump file is required run `docker exec -t map-search-service-database pg_dump -c -U postgres -d renti_db > ~/Downloads/renti-dump.sql`. The file will be saved in your `Downloads` folder

## Using the endpoints

The server runs on port 8080. Example query `http://localhost:8080/api/v1/find-properties?type=rent&suburb=belmore-nsw-2192`

This will:

1. Grab all suburbs within the search radius from the `aus_suburbs` table
2. Check each suburb for existing/unexpired listing data from the cache
3. If the listings exist/are unexpired they will be retireved from the database OR if the listings do not exist they will be scraped, formatted and saved in the database with a record saved in the cache
4. Response is sent to FE

Suburbs (`belmore-nsw-2192` in this example) can be substituted for any suburb Australia wide using the `suburb-state-postcode` format.

Filters can be applied for type of search (rent or buy), bedrooms, bathrooms and parkingby putting a query parameter with the format `&parameter=any-number`/`&parameter=number-any`/`&parameter=number-number`.

Example: `http://localhost:3001/api/v1/find-properties?type=rent&suburb=belmore-nsw-2192&bedrooms=2-5`

## Current bugs

-   Suburbs that have no listings are never added to the cache and hence are always attempted to be scraped
    -   May need to create a timer that removes from cache after X time to try again
    -   Once stored in the database, geolocation data is encoded

## OLD INFO

The back end uses a PostGIS image to support geometry data. Before starting you may choose to restore with a dump file (not required - the scraper will deal with non-exisent data)

1. Create and run image `docker run --name renti-database -p 5438:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=docker -d postgis/postgis`
2. Get container ID `docker ps`
3. docker exec -it <--ID--> psql -U postgres
4. Create database `renti_db`
5. Connect to db (\c renti_db(?)) and `CREATE EXTENSION Postgis;`

-------OLDER--------

1. Create db and port map - `docker run --name renti-db -e POSTGRES_PASSWORD=docker -d -p 5436:5432 postgis/postgis`
2. Build image - `docker build -t renti-db ./`
3. Run container - `docker run -d --name renti-container -p 5436:5432 renti-db`
4. Connect via `psql -h localhost -p 5436 -U postgres -d renti-db`

For PostGIS
https://trevorstanley.medium.com/setup-postgresql-with-postgis-on-docker-8801637a766c
https://registry.hub.docker.com/r/postgis/postgis/ for PostGis
If that fails go here https://dev.to/andre347/

Raw Postgres + use `brew install postgis` > install in db
how-to-easily-create-a-postgres-database-in-docker-4moj
