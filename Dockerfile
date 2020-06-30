FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

# Now copy it into our base image.
FROM scratch

ARG OS=linux
ARG ARCH=amd64

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ADD release/${OS}/${ARCH}/go-wol-server /bin/

CMD ["/bin/go-wol-server"]
