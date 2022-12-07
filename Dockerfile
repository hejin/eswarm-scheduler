FROM debian:stretch-slim

WORKDIR /

COPY bin/eswarm-scheduler /usr/local/bin

CMD ["eswarm-scheduler"]
