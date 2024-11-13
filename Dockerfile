FROM golang:1.21 as Build

WORKDIR /app

COPY . .

RUN make generator loader

FROM mysql:latest as Mysql

COPY --from=Build /app/generator /usr/bin/generator
COPY --from=Build /app/loader /usr/bin/loader