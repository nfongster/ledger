FROM debian:stable-slim

# TODO: Automate this build
COPY cmd/cmd /bin/ledger

CMD ["/bin/ledger"]