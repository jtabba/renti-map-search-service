## Running the server

1. Download `air`
2. `air init`
3. `air` to run the server

## Running the database

The back end uses a PostGIS image to support geometry data. Before starting you may choose to restore with a dump file (not required - the scraper will deal with non-exisent data)

1. Create and run image `docker run --name renti-database -p 5438:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=docker -d postgis/postgis`
2. Get container ID `docker ps`
3. docker exec -it <--ID--> psql -U postgres
4. Create database `renti_db`
5. Connect to db (\c renti_db(?)) and `CREATE EXTENSION Postgis;`

-------OLD--------

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

## Using the endpoints

The server runs on port 3001. Example query `http://localhost:3001/api/v1/find-properties?type=rent&suburb=belmore-nsw-2192`

This will:

1. Grab all suburbs within the search radius from the `suburbs` table
2. Check each suburb for existing/unexpired listing data from the cache
3. If the listings exist/are unexpired then they will be retireved from the database OR if the listings do not exist they will be scraped, formatted and saved in the database with a record saved in the cache 4. Response is sent to FE

Suburbs (`belmore-nsw-2192` in this example) can be substituted for any suburb Australia wide using the `suburb-state-postcode` format.

Filters can be applied for type of search (rent or buy), bedrooms, bathrooms and parkingby putting a query parameter with the format `&parameter=any-number`/`&parameter=number-any`/`&parameter=number-number`.

Example: `http://localhost:3001/api/v1/find-properties?type=rent&suburb=belmore-nsw-2192&bedrooms=2-5`
