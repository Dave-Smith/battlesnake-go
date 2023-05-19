FROM golang:1.20 as build

WORKDIR /usr/src/app

COPY . .
RUN go mod download && go mod verify
RUN go build -v -o /usr/local/bin/app ./...

#FROM alpine:latest
#COPY --from=build /usr/local/bin/app ./
#CMD ["./app"]
CMD ["app"]