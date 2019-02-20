from unittest import mock

import pytest


# Mock RPi
MOCK_RASPY = mock.MagicMock()
MODULES = {
    "RPi": MOCK_RASPY,
    "RPi.GPIO": MOCK_RASPY.GPIO,
}
PATCHER = mock.patch.dict("sys.modules", MODULES)
PATCHER.start()


@pytest.fixture
def mock_rpi():
    return MOCK_RASPY
