FROM golang:1.13-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /gitpages


FROM gcr.io/distroless/base-debian11
WORKDIR /

COPY --from=build /gitpages /gitpages
EXPOSE 2289
USER nonroot:nonroot
ENTRYPOINT ["/gitpages"]
