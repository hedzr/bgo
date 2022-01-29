#
ARG BASE_BUILD_IMAGE
ARG BASE_TARGET_IMAGE

ARG USER_ID
ARG GROUP_ID
ARG USERNAME



#
# Stage I
#
FROM ${BASE_BUILD_IMAGE:-golang:alpine} AS builder

ARG APPNAME
ARG VERSION
ARG PORT

ARG APK_MIRROR="dl-cdn.alpinelinux.org"
ARG BUILDTIME=""
ARG GIT_REVISION=""
ARG W_PKG="github.com/hedzr/cmdr/conf"
ARG GOPROXY="https://goproxy.cn,direct"

ENV APP_HOME="/var/lib/$APPNAME" TGT=/app \
    USER=${USERNAME:-appuser} \
    UID=${USER_ID:-500} GID=${GROUP_ID:-500}

# Install git.
# Git is required for fetching the dependencies.
RUN echo "${APK_MIRROR}"; echo "${APPNAME}"; echo "${APP_HOME}"; sed -i "s/dl-cdn.alpinelinux.org/${APK_MIRROR}/g" /etc/apk/repositories; \
    apk update \
    && apk add --no-cache git ca-certificates tzdata musl-dev musl-utils strace \
    && update-ca-certificates
# RUN apk info -vv | sort

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN mkdir -p $APP_HOME && chown -R ${USER}: $APP_HOME
WORKDIR $APP_HOME
COPY . .
RUN ls -la . && pwd
RUN mkdir -p $TGT/var/lib/$APPNAME/conf.d \
    && mkdir -p $TGT/var/run/$APPNAME $TGT/var/log/$APPNAME $TGT/etc $TGT/etc/sysconfig $TGT/etc/default \
    && touch $TGT/var/lib/$APPNAME/$APPNAME.yml $TGT/etc/sysconfig/$APPNAME $TGT/etc/default/$APPNAME \
    && cp -R ./public $TGT/var/lib/$APPNAME/ \
    && cp -R ./ci/etc/$APPNAME $TGT/etc/ \
    && cp -R ./ci/certs $TGT/etc/$APPNAME/ \
    && chown -R ${USER}: $TGT/var/lib/$APPNAME $TGT/var/log/$APPNAME $TGT/var/run/$APPNAME $TGT/etc/$APPNAME $TGT/etc/sysconfig/$APPNAME $TGT/etc/default/$APPNAME $TGT/var/lib/$APPNAME/$APPNAME.yml
# && touch /target/$APPNAME/var/lib/conf.d/90.alternative.yml
# RUN ls -la ./ &&
# RUN ls -la $TGT/etc/$APPNAME


ENV GOPROXY="$GOPROXY"
RUN echo "Using GOPROXY=$GOPROXY" \
    && go mod download
RUN export GOVER=$(go version) \
    && export VERSION="$(grep -E 'Version[ \t]+=[ \t]+' ./cli/app/doc.go|grep -Eo '[0-9.]+')" \
    && export LDFLAGS="-s -w \
        	-X \"$W_PKG.Buildstamp=$BUILDTIME\" -X \"$W_PKG.Githash=$GIT_REVISION\" \
        	-X \"$W_PKG.Version=$VERSION\" -X \"$W_PKG.GoVersion=$GOVER\" " \
    && echo "Using APPNAME=$APPNAME VERSION=$VERSION" \
    && CGO_ENABLED=0 go build -v -tags docker -tags k8s,istio -tags cmdr-apps \
       -ldflags "$LDFLAGS" \
       -o $TGT/var/lib/$APPNAME/$APPNAME ./cli/your-starter/
RUN ls -la $TGT $TGT/var/lib/$APPNAME $TGT/etc/$APPNAME
# RUN ldd --help
# RUN ldd $TGT/var/lib/$APPNAME/$APPNAME   # need musl-utils & musl-dev
# RUN strace $TGT/var/lib/$APPNAME/$APPNAME  # need strace
RUN rm /var/cache/apk/*
RUN ls -la /var /usr /



#
# Stage II
#
# 1. my-golang-alpine
# FROM ${BASE_TARGET_IMAGE:-scratch}
FROM scratch

ARG APPNAME
ARG VERSION
ARG PORT

LABEL com.hedzr.image.authors="hedzr <hedzrz@gmail.com>"
LABEL com.hedzr.image.description="microservice docker image with hedzr/cmdr"
LABEL description="microservice docker image with hedzr/cmdr"
LABEL version="$VERSION"

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

ENV WF_WFS_CORE_INF_DB_DRIVER=""
ENV WF_WFS_CORE_INF_DB_MYSQL_DSN=""
ENV WF_WFS_CORE_INF_CACHE_DEVEL_PEERS=""
ENV WF_WFS_CORE_INF_CACHE_DEVEL_USER=""
ENV WF_WFS_CORE_INF_CACHE_DEVEL_PASS=""
ENV WF_WFS_CORE_INF_PROD_DEVEL_PEERS=""
ENV WF_WFS_CORE_INF_PROD_DEVEL_USER=""
ENV WF_WFS_CORE_INF_PROD_DEVEL_PASS=""
ENV WF_WFS_CORE_INF_DOCKER_DEVEL_PEERS=""
ENV WF_WFS_CORE_INF_DOCKER_DEVEL_USER=""
ENV WF_WFS_CORE_INF_DOCKER_DEVEL_PASS=""
ENV APP_HOME="/var/lib/$APPNAME" TGT=/app \
    USER="${USERNAME:-appuser}" \
    UID="${USER_ID:-500}" GID="${GROUP_ID:-500}"

WORKDIR "${APP_HOME}"

VOLUME [ "/var/log/$APPNAME", "/var/run/$APPNAME", "/var/lib/$APPNAME/conf.d" ]

# COPY --from=builder /usr/share/i18n /usr/share/i18n
# COPY --from=builder /var/lib/locales /var/lib/locales
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=$USER $TGT/var/lib/$APPNAME /var/lib/$APPNAME
COPY --from=builder --chown=$USER $TGT/var/log/$APPNAME /var/log/$APPNAME
COPY --from=builder --chown=$USER $TGT/var/run/$APPNAME /var/run/$APPNAME
COPY --from=builder --chown=$USER $TGT/etc/$APPNAME /etc/$APPNAME
COPY --from=builder --chown=$USER $TGT/etc/sysconfig/$APPNAME /etc/sysconfig/$APPNAME
#COPY --from=builder --chown=$USER $TGT/etc/default/$APPNAME /etc/default/$APPNAME
EXPOSE $PORT

# Use an unprivileged user.
USER $USER
ENTRYPOINT ["/var/lib/your-starter/your-starter"]
#ENTRYPOINT ["$APP_HOME/$APPNAME"]
CMD ["--help"]
#CMD ["server", "run"]
