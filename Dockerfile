FROM golang:1.11-alpine AS BUILD

RUN apk add --no-cache gcc build-base git mercurial 

ENV BUILD_PATH=$GOPATH/src/github.com/abilioesteves/metrics-generator-tabajara/

RUN mkdir -p ${BUILD_PATH}

WORKDIR ${BUILD_PATH}

ADD ./src ./

RUN go get -v ./...

WORKDIR ${BUILD_PATH}/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /tabajara .

FROM scratch

EXPOSE 3000

COPY --from=BUILD /tabajara /
ADD startup.sh /

CMD [ "/startup.sh" ]
