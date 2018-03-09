FROM golang:alpine as build
WORKDIR /go/src/github.com/blueimp/passphrase
COPY . .
# ldflags explanation (see `go tool link`):
#   -s  disable symbol table
#   -w  disable DWARF generation
RUN go build -ldflags="-s -w" -o /bin/passphrase ./passphrase

FROM scratch
COPY --from=build /bin/passphrase /bin/
USER 65534
ENTRYPOINT ["passphrase"]
