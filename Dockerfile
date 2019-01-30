FROM golang:1.11-alpine AS BUILD

RUN apk add --no-cache gcc build-base git mercurial 

ENV BUILD_PATH=$GOPATH/src/github.com/abilioesteves/metrics-generator-tabajara/src

RUN mkdir -p ${BUILD_PATH}

WORKDIR ${BUILD_PATH}

ADD ./src ./

RUN go get -v ./...

WORKDIR ${BUILD_PATH}/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /tabajara .

FROM alpine:latest

EXPOSE 32856

COPY --from=BUILD /tabajara /

RUN echo "Starting the almighty Metrics Generator Tabajara..."

CMD [ "/tabajara" ]
