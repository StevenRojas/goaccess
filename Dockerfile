FROM golang:alpine AS goaccess

# Working directory for build
WORKDIR /build

# Download dependecies from go.mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build the code to get the binary (goaccess)
COPY . .
WORKDIR /build/cmd/goaccess/
RUN go get
RUN go build -o goaccess

# # Working directory for place the binary
WORKDIR /bin
RUN mkdir init

RUN cp /build/cmd/goaccess/goaccess .
RUN cp -R /build/init/modules/ ./init
EXPOSE 8077
ENTRYPOINT ["/bin/goaccess"]
# RUN ping google.com