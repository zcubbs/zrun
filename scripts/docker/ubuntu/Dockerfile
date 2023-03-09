FROM ubuntu:20.04
MAINTAINER "zcubbs"

# timezone
RUN apt update && apt install -y tzdata dos2unix; \
    apt clean;

# install go
RUN apt update && apt install -y wget; \
    wget https://go.dev/dl/go1.20.2.linux-amd64.tar.gz; \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz; \
    apt clean;

ENV PATH $PATH:/usr/local/go/bin
RUN go version


# sshd
RUN mkdir /run/sshd; \
    apt install -y openssh-server; \
    sed -i 's/^#\(PermitRootLogin\) .*/\1 yes/' /etc/ssh/sshd_config; \
    sed -i 's/^\(UsePAM yes\)/# \1/' /etc/ssh/sshd_config; \
    apt clean;

# entrypoint
RUN { \
    echo '#!/bin/bash -eu'; \
    echo 'ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime'; \
    echo 'echo "root:${ROOT_PASSWORD}" | chpasswd'; \
    echo 'exec "$@"'; \
    } > /usr/local/bin/entry_point.sh; \
    chmod +x /usr/local/bin/entry_point.sh;

ENV TZ Europe/Paris

ENV ROOT_PASSWORD root

EXPOSE 22

ENTRYPOINT ["tail", "-f", "/dev/null"]
CMD    ["/usr/sbin/sshd", "-D", "-e"]
