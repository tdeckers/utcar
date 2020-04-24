FROM scratch

ADD utcar /
ADD ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/utcar"]