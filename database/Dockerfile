FROM postgres:10.3

# Copy file to run on image deploy
COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]