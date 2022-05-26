FROM golang:1.16-alpine as build
WORKDIR /build

ARG VERSION
ARG COMMIT_HASH
ARG BUILD_DATE

# add some necessary packages
RUN apk update && \
    apk add libc-dev && \
    apk add gcc && \
    apk add make

# prevent the re-installation of vendors at every change in the source code
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# Copy and build the app
COPY . .
RUN echo "[version]: ${VERSION} | [commit_hash]: ${COMMIT_HASH} | [build_date]: ${BUILD_DATE}"

RUN export VERSION_TO_ENV="${VERSION}" \
    && export COMMIT_HASH_TO_ENV="${COMMIT_HASH}" \
    && export BUILD_DATE_TO_ENV="${BUILD_DATE}" \
    && go build -ldflags "-X main.version=${VERSION_TO_ENV} -X main.commitHash=${COMMIT_HASH_TO_ENV}" -o api ./cmd/serverd
#RUN go build -o api ./cmd/serverd

FROM golang:1.16-alpine as dist
WORKDIR /app

COPY --from=build /build/api /app/api

LABEL org.opencontainers.image.source="https://github.com/VIIGstar/scaffold-api-server"

RUN mkdir config

CMD mv ./config/conf.toml . && ./api