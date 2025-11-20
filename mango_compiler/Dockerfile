# commands to build and run docker image:
# docker build -t mango .
# docker run --rm -it -v $(pwd):/usr/src/app mango

# use docker to create isolated ubutu environment for mango compiler
FROM ubuntu:24.04

# prevent interactive prompts
ENV DEBIAN_FRONTEND=noninteractive

# install essential build tools
RUN apt-get update && \
    apt-get install -y \
        build-essential \
        wget \
        cmake \
        flex \
        bison \
        git \
        g++ \
        clang \
        ninja-build \
        pkg-config \
        libstdc++-12-dev \
        libfl-dev \
        && rm -rf /var/lib/apt/lists/*

# install latest bison version
RUN wget https://ftp.gnu.org/gnu/bison/bison-3.8.2.tar.gz && \
    tar -xzf bison-3.8.2.tar.gz && \
    cd bison-3.8.2 && \
    ./configure --prefix=/usr/local && \
    make -j$(nproc) && make install && \
    cd .. && rm -rf bison-3.8.2 bison-3.8.2.tar.gz

# set working directory
WORKDIR /usr/src/app

# copy source code from host
COPY . .

# start bash shell in container
CMD ["bash"]
