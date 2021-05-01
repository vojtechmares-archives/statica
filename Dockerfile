FROM ubuntu:20.04

COPY statica /

ENTRYPOINT [ "/statica" ]
