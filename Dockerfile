FROM golang:1.11-rc-alpine as build
WORKDIR /opt/passphrase
COPY . .
# ldflags explanation (see `go tool link`):
#   -s  disable symbol table
#   -w  disable DWARF generation
RUN cd ./passphrase && go build -ldflags="-s -w" -o /bin/passphrase

FROM scratch
COPY --from=build /bin/passphrase /bin/
USER 65534
ENTRYPOINT ["passphrase"]
