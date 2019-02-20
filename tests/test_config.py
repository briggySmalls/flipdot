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
    parser = ConfigParser.create(config_file_path)

    # Assert
    assert parser.basic_config['serial_port'] == '/dev/ttyUSB0'
    assert parser.basic_config['grpc_port'] == 5001
    assert parser.pin_config.sign == 40
    assert parser.pin_config.light == 38
    # Assert first sign
    assert parser.signs_config[0].name == 'top'
    assert parser.signs_config[0].address == 1
    assert parser.signs_config[0].width == 84
    assert parser.signs_config[0].height == 7
    assert parser.signs_config[0].flip
    # Assert second sign
    assert parser.signs_config[1].name == 'bottom'
    assert parser.signs_config[1].address == 2
    assert parser.signs_config[1].width == 12
    assert parser.signs_config[1].height == 18
    assert not parser.signs_config[1].flip
