FROM golang:1.15.4-alpine

WORKDIR /src

COPY . .

CMD ["/bin/cfc_suggestions"]

