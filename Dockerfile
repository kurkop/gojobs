FROM golang:1.13 AS build

WORKDIR /go/src/app
ENV CGO_ENABLED=0
ENV IN_CLUSTER=true

COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gojobs-api cmd/gojobs-api/main.go

FROM scratch AS bin
ENV IN_CLUSTER=true
COPY --from=build /go/src/app/gojobs-api /
EXPOSE 8000
CMD ["/gojobs-api"]