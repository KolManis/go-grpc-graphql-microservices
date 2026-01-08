FROM postgres:15-alpine

COPY up.sql /docker-entrypoint-initdb.d/

CMD ["postgres"]