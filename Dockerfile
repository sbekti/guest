FROM node:16 AS node-builder
WORKDIR /build

COPY package.json yarn.lock ./
COPY public/ ./public
COPY src/ ./src

RUN yarn install --network-timeout 1000000 && yarn build

FROM golang:latest AS go-builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build

# Let's cache modules retrieval - those don't change so often
# COPY go.mod .
# COPY go.sum .
RUN go mod init github.com/sbekti/guest
RUN go mod download

# Copy the code necessary to build the application
# You may want to change this to copy only what you actually need.
COPY . .

# Build the application
RUN go build -v -o main

# Let's create a /dist folder containing just the files necessary for runtime.
# Later, it will be copied as the / (root) of the output image.
WORKDIR /dist
RUN cp /build/main ./main

# Optional: in case your application uses dynamic linking (often the case with CGO), 
# this will collect dependent libraries so they're later copied to the final image
# NOTE: make sure you honor the license terms of the libraries you copy and distribute
# RUN ldd main | tr -s '[:blank:]' '\n' | grep '^/' | \
#     xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
# RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/

# Copy or create other directories/files your app needs during runtime.
# E.g. this example uses /data as a working directory that would probably
#      be bound to a perstistent dir when running the container normally
RUN mkdir /data

# Create the minimal runtime image
FROM scratch

COPY --chown=0:0 --from=go-builder /dist /
COPY --chown=0:0 --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set up the app to run as a non-root user inside the /data folder
# User ID 65534 is usually user 'nobody'. 
# The executor of this image should still specify a user during setup.
COPY --chown=65534:0 --from=go-builder /data/ /data
COPY --chown=65534:0 --from=node-builder /build/build/ /data/build
USER 65534
WORKDIR /data

ENTRYPOINT ["/main"]
