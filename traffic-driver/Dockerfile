FROM fedora

ADD entrypoint.sh entrypoint.sh

RUN curl https://hey-release.s3.us-east-2.amazonaws.com/hey_linux_amd64 -o hey
RUN chmod +x hey

RUN ls
CMD ["./entrypoint.sh"]
