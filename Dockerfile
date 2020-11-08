FROM golang:1.15.4-alpine

WORKDIR /src

COPY . .

RUN apk add --no-cache gcc musl-dev
RUN go build -o /bin/cfc_suggestions .

CMD ["/bin/cfc_suggestions"]

