FROM alpine:3.13.5
LABEL version="1.0.0"
LABEL alpineversion="3.13.5"
LABEL golangversion="1.16.4"
LABEL homepage="https://github.com/hekmon/vigixporter"
ENV VIGIXPORTER_STATIONS=""
ENV VIGIXPORTER_VMURL=""
ENV VIGIXPORTER_VMUSER=""
ENV VIGIXPORTER_VMPASS=""
COPY vigixporter_alpine /usr/local/bin/vigixporter
VOLUME /var/lib/vigixporter
WORKDIR /var/lib/vigixporter
ENTRYPOINT /usr/local/bin/vigixporter