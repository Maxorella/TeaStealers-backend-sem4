FROM postgres:latest

RUN apt-get update && apt-get install -y postgresql-contrib postgis

RUN echo "CREATE EXTENSION IF NOT EXISTS \"postgis\";" >> /docker-entrypoint-initdb.d/init.sql

# psql -U maxorella -d maxorella_db
# \dt

# FROM migrate/migrate:v4.15.2
# COPY ./db/migrations /migrations