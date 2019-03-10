#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""The setup script."""

from setuptools import setup

with open('README.rst') as readme_file:
    readme = readme_file.read()

with open('HISTORY.rst') as history_file:
    history = history_file.read()

requirements = [
    'Click>=6.0',
    'pyflipdot>=0.1.3',
    'pyserial>=3.4',
    'RPi.GPIO>=0.6.5',
    'grpcio>=1.9.0',
    'numpy==1.9.3',
    'protobuf>=3.6.0',
    'toml>=0.10.0',
    'grpcio-reflection>=1.19.0',
]

setup_requirements = ['pytest-runner', ]

test_requirements = ['pytest', ]

setup(
    author="Sam Briggs",
    author_email='briggySmalls90@gmail.com',
    classifiers=[
        'Development Status :: 2 - Pre-Alpha',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Natural Language :: English',
        "Programming Language :: Python :: 2",
        'Programming Language :: Python :: 2.7',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.4',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
    ],
    description="Controller for Hanover flipdot signs",
    entry_points={
        'console_scripts': [
            'flipdot_controller=flipdot_controller.cli:main',
        ],
    },
    install_requires=requirements,
    license="MIT license",
    long_description=readme + '\n\n' + history,
    include_package_data=True,
    keywords='flipdot_controller',
    name='flipdot_controller',
    packages=['flipdot_controller', 'flipdot_controller.protos'],
    setup_requires=setup_requirements,
    test_suite='tests',
    tests_require=test_requirements,
    url='https://github.com/briggySmalls/flipdot_controller',
    version='0.1.0',
    zip_safe=False,
)
