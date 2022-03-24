FROM golang:1.13-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go get -v -t -d ./...

COPY . ./
RUN go build -o /gitpages


FROM gcr.io/distroless/base-debian10
WORKDIR /

COPY --from=build /gitpages /gitpages
EXPOSE 2289
USER nonroot:nonroot
ENTRYPOINT ["/gitpages"]
