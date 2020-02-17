#FROM hashicorp/terraform:0.11.14 as terraform
FROM golang:1.13-alpine3.10

ARG PROJECT
ARG USER_ID
ARG AWS_DEFAULT_REGION
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG VAULT_TOKEN

ENV BUILD_DIR /builder/src/aws-api
ENV BUILD_SCRIPT "$BUILD_DIR"/docker/build.sh

COPY . "$BUILD_DIR"
WORKDIR "$BUILD_DIR"
VOLUME "$BUILD_DIR"

#COPY --from=terraform /bin/terraform /bin/terraform

# installs
RUN apk update \
  && apk add -u sed vim curl git jq bash openssh-client


RUN chmod +x "$BUILD_SCRIPT" \
 && chmod +x "$BUILD_DIR"/scripts/* \
 && adduser -h /builder/ -D -u $USER_ID docker \
 && chown docker:users -R /builder

USER docker

ENTRYPOINT "$BUILD_SCRIPT"