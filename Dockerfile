FROM scratch

COPY statica /

ENTRYPOINT [ "/statica" ]
