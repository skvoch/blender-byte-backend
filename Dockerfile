FROM golang:1.12 as build
WORKDIR /go/src/app
COPY . .
ENV GO111MODULE on
RUN go get ./...
RUN make
# Now copy it into our base image.
FROM gliderlabs/alpine:3.4
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add --no-cache ca-certificates 
RUN apk add --no-cache bash
RUN apk add --no-cache python3
COPY --from=build /go/src/app/back /back
COPY --from=build /go/src/app/python/* /python/
COPY --from=build /go/src/app/private.json /private.json
COPY --from=build /go/src/app/start.sh /start.sh
RUN pip3 install --upgrade pip setuptools
RUN pip3 install -r /python/requirements.txt

ENTRYPOINT ["/bin/bash", "/start.sh"]
