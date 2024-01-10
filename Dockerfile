FROM scratch
WORKDIR /
COPY m3u8-downloader /
USER 65534
ENTRYPOINT ["/m3u8-downloader"]