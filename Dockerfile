FROM golang:1.20.4

ENV WORK_DIR /opt/env-leaker

RUN mkdir -p $WORK_DIR
WORKDIR $WORK_DIR

COPY env-leaker $WORK_DIR/env-leaker

CMD [ "/opt/env-leaker/env-leaker", "-o", "stdout" ]

