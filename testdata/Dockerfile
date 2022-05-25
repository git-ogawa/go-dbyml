from alpine:latest

label "builder" go-dev

FROM alpine:latest
ARG key1
RUN echo "$key1" > arg.txt && \
    cat arg.txt && \
    rm arg.txt

run echo aaa