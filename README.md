# Jim's-magic-sign

[![CircleCI](https://circleci.com/gh/briggySmalls/flipdot.svg?style=svg)](https://circleci.com/gh/briggySmalls/flipdot)
![GitHub](https://img.shields.io/github/license/briggySmalls/flipdot.svg)

A collection of applications that work together to make a smart flipdot display

![](https://imgur.com/hQYodeh.gif)

## Architecture

The following deployment diagram illustrates the different components of the flipdot project, and how they are deployed. Beware, this diagram is probably trying to say too much!

![](http://www.plantuml.com/plantuml/png/VL91I_D04BtFhvZZznqYHIf82A6sWdYoMDH3yR19ndH9cbrcDzOW_Uzkih5kRN4E2Vioy-QzjvaPAzYssnfC9HijM6pH0V9Dv1O_0Lrb8gzALcrJB5Ij69TgLn3FwvREVKuIkv5SeAEoNPhYoqPQMcrLHR07Q5H1oCBeZ9WxBSSLJBaLJLaJ5YglY7juh8COeJMk0ODAP5egk71r36Ufwpr0ht3ALR1y9pwbqvtTgEOifH_varMp-kZmTm37Iyh7vIBQKUQR0xh-kVUalTFetoGQPSR3K8otm-cdO_8_yYpV3JVEjQC8m-nV0S1KYouuPwsrpY-CUKHHEix4-BIQ1x2VZF5kUVy0qAzF7EZ7HpIWDR9ip7ZP6QkTyJSQPjL7i8OWUjjfCZbBRNtyGVIk5tn8pjtWs9jtU7m8rDltwdl5NaICRABfzHKC-aFySleCCseAz-j_0G00)

Software items are denoted with the 📄 symbol, and those that are blue can be found within this repo:

item | description
--- | ---
[app](./app) | Application that displays messages and the time on the signs
[driver](./driver) | Driver responsible for sending commands to signs
[web](./web) | Website that collects messages for displaying on the signs
[protos](./protos) | Protocol buffers that define interface between software items (marked in diagram as )
[proxy](https://github.com/improbable-eng/grpc-web/tree/master/go/grpcwebproxy) | reverse proxy from [improbable-eng](https://github.com/improbable-eng) allowing for the gRPC services to be consumed from browsers.
