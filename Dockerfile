ARG work_dir=/go/src/github.com/Bnei-Baruch/gxydb-api
ARG build_number=dev

FROM golang:1.21 AS build

LABEL maintainer="edoshor@gmail.com"

ARG work_dir
ARG build_number

ENV GOOS=linux \
	CGO_ENABLED=0

WORKDIR ${work_dir}
COPY . .
RUN go build -ldflags "-w -X github.com/Bnei-Baruch/gxydb-api/version.PreRelease=${build_number}"

FROM alpine:latest

ARG work_dir
WORKDIR /app
COPY ./misc/wait-for /wait-for
COPY --from=build ${work_dir}/gxydb-api .

EXPOSE 8080
CMD ["./gxydb-api", "server"]
