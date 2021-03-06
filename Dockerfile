FROM golang:alpine

# Install go dep
RUN apk add --no-cache curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN apk del curl

# Init dep project
RUN mkdir -p $GOPATH/src/github.com/jessebmiller/trbac
WORKDIR $GOPATH/src/github.com/jessebmiller/trbac

# Install dependencies
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["trbac"]