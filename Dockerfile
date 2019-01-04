FROM golang:1.11-alpine AS BUILD

RUN apk add --no-cache gcc build-base git mercurial 

ENV BUILD_PATH=$GOPATH/src/github.com/abilioesteves/metrics-generator-tabajara/

RUN mkdir -p ${BUILD_PATH}

WORKDIR ${BUILD_PATH}

ADD ./src ./

RUN go get -v ./...

WORKDIR ${BUILD_PATH}/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /metrics-generator-tabajara .

FROM scratch

EXPOSE 3000

ENV SERVER_NAME ''
ENV COMPONENT_NAME 'testserver'
ENV COMPONENT_VERSION '1.0.0'
ENV ACCIDENT_RESOURCE ''
ENV ACCIDENT_TYPE ''
ENV ACCIDENT_RATIO 1

COPY --from=BUILD /metrics-generator-tabajara /
ADD startup.sh /

CMD [ "/startup.sh" ]
