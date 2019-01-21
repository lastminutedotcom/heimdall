ARG TAG_NAME
FROM registry.bravofly.intra:5000/application/goappfw:$TAG_NAME

LABEL author="Team SRE <sre@lastminute.com>"

ARG APP_NAME

COPY $APP_NAME /application

ENV APPFW_NAME $APP_NAME

ENTRYPOINT ["/run_helper.sh"]