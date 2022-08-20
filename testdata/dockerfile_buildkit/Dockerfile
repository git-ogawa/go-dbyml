# syntax=docker/dockerfile:1.2
FROM alpine:3.16.0

ARG key1
RUN echo "$key1" > arg.txt && \
    cat arg.txt && \
    rm arg.txt

run echo aaa
