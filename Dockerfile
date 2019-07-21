# 1) build environment
FROM ubuntu:18.04 as builder
ENV GOLANG_VERSION 1.11.9

RUN apt-get update && apt-get install -y wget git gcc curl musl-dev libtool

# install golang
RUN wget -P /tmp https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf /tmp/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    rm /tmp/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# install dep
WORKDIR ${GOPATH}/src/github.com/mariusgiger/iot-abe
RUN curl -s https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY Gopkg.toml Gopkg.lock ./

# install dependencies and fix eth bug
RUN go get -v "github.com/ethereum/go-ethereum/crypto/secp256k1"  && \
    dep ensure -v -vendor-only && \
    cp -r "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"

# install c crypto dependencies
RUN apt-get install -y libglib2.0-dev flex bison libssl-dev libgmp-dev
RUN wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz -O ./pbc-0.5.14.tar.gz && \
    tar -zxvf pbc-0.5.14.tar.gz && cd pbc-0.5.14/ && \
    ./configure && \
    make && \
    make install
RUN wget http://hms.isi.jhu.edu/acsc/cpabe/libbswabe-0.9.tar.gz -O ./libbswabe-0.9.tar.gz && \
    tar -zxvf libbswabe-0.9.tar.gz && \
    cd libbswabe-0.9 && \
    ./configure && \ 
    make && \ 
    make install

COPY pkg/crypto/policy_lang.y.linux ./policy_lang.y
COPY pkg/crypto/Makefile.linux Makefile
RUN wget http://hms.isi.jhu.edu/acsc/cpabe/cpabe-0.11.tar.gz -O ./cpabe-0.11.tar.gz && \
    tar -zxvf cpabe-0.11.tar.gz && \
    cd cpabe-0.11 && \ 
    ./configure --with-gmp-lib=/usr/lib/x86_64-linux-gnu && \ 
    #fixes missing semicolon 
    cp ../policy_lang.y policy_lang.y && \
    #fixes adding -lgmp to linker command (At THE END!) 
    cp ../Makefile Makefile && \
    make && \ 
    make install

# build app
WORKDIR ${GOPATH}/src/github.com/mariusgiger/iot-abe
COPY . ./
RUN export VERSION=$(git describe --tags 2>/dev/null) && \
    export GITHASH=$(git log -1 --format='%H') && \
    export BUILDTIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ') && \
    GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -tags netgo -a -o /out/iot-abe \
    -ldflags "-X github.com/mariusgiger/iot-abe/cmd.Version=${VERSION} -X github.com/mariusgiger/iot-abe/cmd.GitHash=${GITHASH} -X github.com/mariusgiger/iot-abe/cmd.BuildTime=${BUILDTIME} -w -extldflags '-static'" .

# 2) deploy environment
FROM alpine:latest as production
WORKDIR /app
RUN apk --no-cache add ca-certificates libc6-compat
COPY --from=builder /out/iot-abe /app/bin/iot-abe
COPY config.yml ./config.yml
RUN /bin/sh -c "bin/iot-abe -c './config.yml' --help"
EXPOSE 8080
ENTRYPOINT ["bin/iot-abe"]
CMD ["server"]