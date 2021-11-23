To be able to cache build artifact in docker repository we have a separated docker builder image. To build the application image:

    docker build -t builder -f docker/Dockerfile.builder .
