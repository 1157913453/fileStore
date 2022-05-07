FROM debian

ENV REDIS_VER 6.2.6
ENV TARBALL http://download.redis.io/releases/redis-$REDIS_VER.tar.gz

RUN echo "==> 安装'curl','make','gcc'..." && \
    apt-get update && \
    apt-get install -y curl gcc make && \
    \
    echo "==> 下载安装Redis" && \
    curl -L $TARBALL | tar zxv && \
    cd redis-$REDIS_VER && \
    make && \
    make install && \
    rm -rf /var/lib/apt/lists/* /redis-$REDIS_VER


CMD ["redis-server"]