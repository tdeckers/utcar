FROM scratch

ADD utcar /
ADD ca-certificates.crt /etc/ssl/certs/

EXPOSE 12300

ENTRYPOINT ["/utcar"]