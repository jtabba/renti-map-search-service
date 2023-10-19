FROM postgis/postgis
ENV POSTGRES_PASSWORD docker
# ENV POSTGRES_DB renti_db
# COPY map_search_listings.sql /docker-entrypoint-initdb.d/


# COMMANDS
# view db: psql -h localhost -p 5436 -U postgres -d map_search_listings