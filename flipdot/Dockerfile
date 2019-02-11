FROM python:3.7-stretch

# Prepare package (generate source)
RUN pip install invoke grpcio-tools
COPY . /package/
WORKDIR /package/
# Build package into a distribution
RUN invoke dist

FROM balenalib/raspberrypi3:stretch

RUN apt-get update
# Note: We install python ourselves as docker doesn't compile CPython with fpectl
# If we took Docker's we wouldn't be able to use piwheels
RUN apt-get install libatlas3-base python3 python3-pip python3-setuptools

# Make piwheels available
RUN printf "[global]\nextra-index-url=https://www.piwheels.org/simple\n" > /etc/pip.conf

# Ensure we can compile an efficient package installation
RUN pip3 install wheel

# Copy build package into a new image
COPY --from=0 /package/dist/flipdot_controller*.tar.gz /app/
WORKDIR /app/

# Install the package
RUN pip3 install flipdot_controller*.tar.gz

ENV GRPC_PORT=5001
CMD flipdot_controller --serial-port $SERIAL_PORT --grpc-port $GRPC_PORT
