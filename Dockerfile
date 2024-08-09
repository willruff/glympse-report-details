# syntax=docker/dockerfile:experimental
FROM python:3.12-slim
WORKDIR /app

ARG RPT_DATE
ARG BEGIN_DATE
ARG END_DATE
ARG ORG_DOWNLOAD_ID
ARG ORG_UPLOAD_ID
ARG BUCKET_NAME
ARG COLUMN_SELECT
ARG SEARCH_FILTER

ENV XSV_VERSION=0.13.0
ENV XSV_SHA256=271e798160472830d7151673383afaba4c37209673f5157cf37e8f5e308f1cac
ENV RPT_DATE=${RPT_DATE}
ENV BEGIN_DATE=${BEGIN_DATE}
ENV END_DATE=${END_DATE}
ENV ORG_DOWNLOAD_ID=${ORG_DOWNLOAD_ID}
ENV ORG_UPLOAD_ID=${ORG_UPLOAD_ID}
ENV BUCKET_NAME=${BUCKET_NAME}
ENV COLUMN_SELECT=${COLUMN_SELECT}
ENV SEARCH_FILTER=${SEARCH_FILTER}

RUN apt-get update && apt-get install -y \
      curl \
      jq \
      python3 \
      python3-pip \
      python3-six && \
    curl -qsL https://github.com/BurntSushi/xsv/releases/download/${XSV_VERSION}/xsv-${XSV_VERSION}-x86_64-unknown-linux-musl.tar.gz | tar xzvf - -C /usr/bin xsv && \
    [ "${XSV_SHA256}  /usr/bin/xsv" = "$(sha256sum /usr/bin/xsv)" ] && \
    pip3 install awscli && \
    apt-get -y remove python3-pip && \
    apt-get -y autoremove && \
    rm -rf /root/.cache

COPY requirements.txt .
RUN --mount=type=ssh pip install -r requirements.txt 

COPY . /app

CMD ["python", "main.py"]
