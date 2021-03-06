# Jim's-magic-sign ~ driver

Driver 'service' for controlling flipdot signs connected to a raspberry pi, written in Python.

## Features

- Sends Hanover flipdot information signals over serial port
- Controls power out using configurable GPIO pins
- Exposes driver gRPC 'service' for network control

## Installation

All development dependencies are tracked using pipenv, thus may be installed with:

```bash
# Install development dependencies
pipenv install --dev

# Start a shell within virtual environment
pipenv shell
```

To see what development tasks may be run, use invoke:

```bash
# From within pipenv virtual environment
invoke --list
```

Before being able to run tests or install the driver, you will need to generate gRPC sources:

```bash
invoke proto
```

You will now have the command line tool available:
```bash
flipdot_controller --help
```

## Configuration

The application is totally configured using a TOML formatted configuration file.

See [config](./config.toml) for a complete example of configurable parameters.

## Docker

A docker container, ready to run on a raspberry pi (ARM7), can be built in the usual way:

```bash
docker build -t driver .
```

Note that this container will need access to both a serial port and the GPIOs:

```bash
docker run \
    -p 5001:5001 \
    --device /dev/ttyUSB0:/dev/ttyUSB \
    --device /dev/gpiomem:/dev/mem \
    driver \
    --config path/to/custom/config.toml
```

The built version of this container can be found on dockerhub as `briggysmalls/flipdot_controller`.

Credits
-------

This package was created with [Cookiecutter](https://github.com/audreyr/cookiecutter) and the [briggySmalls/cookiecutter-pypackage](https://github.com/briggySmalls/cookiecutter-pypackage) project template.
