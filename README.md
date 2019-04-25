# Jim's-magic-sign

[![CircleCI](https://circleci.com/gh/briggySmalls/flipdot.svg?style=svg)](https://circleci.com/gh/briggySmalls/flipdot)
![GitHub](https://img.shields.io/github/license/briggySmalls/flipdot.svg)

A collection of applications that work together to make a smart flipdot display

![](https://imgur.com/hQYodeh.gif)

## Architecture

The following deployment diagram illustrates the different components of the flipdot project, and how they are deployed. Beware, this diagram is probably trying to say too much!

![](http://www.plantuml.com/plantuml/png/VP71JlCm48Jl-nG-_l_1g495Gb4KbHO9SQAY0WuL1oSngPN4jjOE9K9zT-9kA2wjs4EYJ7PdnvzdpWlqNTk0gvMs0aNB1a6zYSApJs13vQAeApITBXUcCSYMPbjAd3UTFFSxJVr6OSciGDzd6NlPA2zNhQabx02qAIL3uMmk4NkhnXKe2ozqrKXMcgAMI7Aedp32MhO65bMgiih0nL37Sfu9Qm_IAvnwbQZU9PxQsTvlZ3vhIID_kbeq7ptx3U1uIVMuNF2jpAaviWlF7Pc6a_9_4px6JFP32Rkb18UEoNzEBjyDDzO6nZ7Cy8e4b8tee--yyzveW947vzauax2drJoIQJ9XTylx1u2mFXr46YSrtkjKfTrusJcQhQCRZLa516k8qFlUIbUWifxmH-Y7NY18Em22p3aF5ic19vsUmeV0b25XwFZq-WhsyMFzSleCCwdBhcs-0000)

Software items are denoted with the 📄 symbol, and those that are blue can be found within this repo:

item | description
--- | ---
[app](./app) | Application that displays messages and the time on the signs
[driver](./driver) | Driver responsible for sending commands to signs
[web](./web) | Website that collects messages for displaying on the signs
[proxy](https://github.com/improbable-eng/grpc-web/tree/master/go/grpcwebproxy) | reverse proxy from [improbable-eng](https://github.com/improbable-eng) allowing for the gRPC services to be consumed from browsers.
