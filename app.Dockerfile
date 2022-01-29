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

ARG APK_MIRROR="dl-cdn.alpinelinux.org"			# ARG APK_MIRROR="mirrors.tuna.tsinghua.edu.cn"
ARG BUILDTIME=""
ARG GIT_REVISION=""
ARG W_PKG="github.com/hedzr/cmdr/conf"
ARG GOPROXY="https://goproxy.cn,direct"

ENV APP_HOME="/var/lib/$APPNAME" TGT=/app \
	USER=${USERNAME:-appuser} \
	UID=${USER_ID:-500} \
	GID=${GROUP_ID:-500}

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
	&& mkdir -p $TGT/var/run/$APPNAME $TGT/var/log/$APPNAME $TGT/etc $TGT/etc/sysconfig $TGT/etc/default $TGT/etc/$APPNAME \
	&& touch $TGT/etc/sysconfig/$APPNAME $TGT/etc/default/$APPNAME \
	&& chown -R ${USER}: $TGT/var/lib/$APPNAME $TGT/var/log/$APPNAME $TGT/var/run/$APPNAME $TGT/etc/$APPNAME $TGT/etc/sysconfig/$APPNAME $TGT/etc/default/$APPNAME
# && touch /target/$APPNAME/var/lib/conf.d/90.alternative.yml
# RUN ls -la ./ &&
# RUN ls -la $TGT/etc/$APPNAME


ENV GOPROXY="$GOPROXY"
RUN echo "Using GOPROXY=$GOPROXY" && go mod download
# 	&& export VERSION="$(grep -E 'Version[ \t]+=[ \t]+' ./cli/bgo/cmdr/doc.go|grep -Eo '[0-9.]+')"
# --mount=type=cache,target=/root/.cache/go-build
RUN export GOVER=$(go version) \
	&& export LDFLAGS="-s -w \
	-X \"$W_PKG.Buildstamp=${BUILDTIME}\" -X \"$W_PKG.Githash=${GIT_REVISION}\" \
	-X \"$W_PKG.Version=${VERSION}\" -X \"$W_PKG.GoVersion=${GOVER}\" \
	-X \"$W_PKG.ServerID=docker-build\" " \
	&& echo "Using APPNAME=$APPNAME VERSION=$VERSION" \
	&& CGO_ENABLED=0 go build -v -trimpath \
	-tags docker -tags k8s,istio -tags cmdr-apps \
	-ldflags "$LDFLAGS" \
	-o $TGT/var/lib/$APPNAME/$APPNAME ./main.go
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
# FROM scratch
FROM golang:alpine

ARG APPNAME
ARG VERSION
ARG PORT

LABEL com.hedzr.image.authors="hedzr <hedzrz@gmail.com>"
LABEL com.hedzr.image.description="make go building easier with hedzr/bgo"
LABEL description="make go building easier with hedzr/bgo"
LABEL version="$VERSION"
LABEL org.opencontainers.image.description="make go building easier with hedzr/bgo"
LABEL org.opencontainers.image.author="hedzr <hedzrz@gmail.com>"
LABEL org.opencontainers.image.url="https://github.com/hedzr/bgo"
LABEL org.opencontainers.image.version="$VERSION"
LABEL org.opencontainers.image.license="Apache 2.0"


# Import the user and group files from the builder.
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /etc/group /etc/group

ENV APP_HOME="/var/lib/$APPNAME" TGT=/app \
	USER="${USERNAME:-appuser}" \
	UID="${USER_ID:-500}" \
	GID="${GROUP_ID:-500}" \
	GOPROXY=$GOPROXY

WORKDIR "/app"
VOLUME ["/app", "/go/pkg", "/tmp"]
# VOLUME [ "/var/log/$APPNAME", "/var/run/$APPNAME", "/var/lib/$APPNAME/conf.d" ]

# COPY --from=builder /usr/share/i18n /usr/share/i18n
# COPY --from=builder /var/lib/locales /var/lib/locales
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder --chown=$USER $TGT/var/lib/$APPNAME /var/lib/$APPNAME
COPY --from=builder $TGT/var/lib/$APPNAME /var/lib/$APPNAME
# COPY --from=builder --chown=$USER $TGT/var/log/$APPNAME /var/log/$APPNAME
# COPY --from=builder --chown=$USER $TGT/var/run/$APPNAME /var/run/$APPNAME
# COPY --from=builder --chown=$USER $TGT/etc/$APPNAME /etc/$APPNAME
# COPY --from=builder --chown=$USER $TGT/etc/sysconfig/$APPNAME /etc/sysconfig/$APPNAME
#COPY --from=builder --chown=$USER $TGT/etc/default/$APPNAME /etc/default/$APPNAME
# EXPOSE $PORT

# Use an unprivileged user.
#USER $USER
ENTRYPOINT ["/var/lib/bgo/bgo"]
#ENTRYPOINT ["$APP_HOME/$APPNAME"]
#CMD ["--help"]
#CMD ["server", "run"]
