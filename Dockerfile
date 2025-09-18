FROM debian:stable-slim

COPY cmd/cmd /bin/ledger

COPY config.json config.json

CMD ["/bin/ledger"]