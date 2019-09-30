FROM golang:1.13-alpine3.10 as builder
ENV app="olist"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
WORKDIR /go/src
ADD ./ ./"${app}"
RUN apk add --no-cache --virtual git gcc g++ musl-dev
RUN (cd $app;go mod tidy;go mod vendor;go build -a -ldflags '-linkmode external -extldflags "-static"' -o ${app} .)

FROM scratch
ENV app="olist"
EXPOSE 80
COPY --from=builder /go/src/"${app}"/"${app}" /"${app}"
ENTRYPOINT [ "/olist"]
CMD [""]
