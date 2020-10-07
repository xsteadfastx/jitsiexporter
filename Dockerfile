FROM scratch
COPY jitsiexporter /
EXPOSE 9700
ENTRYPOINT ["/bin/jitsiexporter", "-debug=true", "-host=0.0.0.0"]
