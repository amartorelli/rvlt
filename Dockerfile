FROM alpine:3.4
COPY helloworld /helloworld

ENTRYPOINT [ "/helloworld" ]
