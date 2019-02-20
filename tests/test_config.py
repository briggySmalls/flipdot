from flipdot_controller.config import ConfigParser


def test_simple_config(tmp_path):
    # Create a dummy config
    config_text = """
        serial_port='/dev/ttyUSB0'
        grpc_port=5001

        [pins]
        sign=40
        light=38

        [signs]
        [signs.top]
        address=1
        width=84
        height=7
        flip=true

        [signs.bottom]
        address=2
        width=12
        height=18
        flip=false
    """
    # Write the config to a file
    config_file_path = tmp_path.joinpath('config.toml')
    with config_file_path.open('w') as file:
        file.write(config_text)
    # Read and validate the config with the parser
    config_parser = ConfigParser.create(config_file_path)
    config = config_parser.config

    # Assert
    assert config['serial_port'] == '/dev/ttyUSB0'
    assert config['grpc_port'] == 5001
    assert config['pins']['sign'] == 40
    assert config['pins']['light'] == 38
    # Assert first sign
    assert config['signs']['top']['address'] == 1
    assert config['signs']['top']['width'] == 84
    assert config['signs']['top']['height'] == 7
    assert config['signs']['top']['flip']
    # Assert second sign
    assert config['signs']['bottom']['address'] == 2
    assert config['signs']['bottom']['width'] == 12
    assert config['signs']['bottom']['height'] == 18
    assert not config['signs']['bottom']['flip']

