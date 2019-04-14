# web

Web application that is a gRPC-web client for sending messages to a flipdot sign, built using:

- Vue.js
- [@improbable-eng/grpc-web](https://github.com/improbable-eng/grpc-web/tree/master/client/grpc-web)

## Project setup
```
yarn install
```

### Development

Start a development server with:

```
yarn run serve
```

The existing `.env.development` configures the app to connect to a gRPC-web server at http://localhost:5002.

### Production

Create a `.env.production` file of the following form:

```ini
# Address where grpc-web server can be contacted
VUE_APP_GRPC_SERVER_ADDRESS=https://address.of.grpc
```

Then build the application:

```
yarn run build
```

### Run your tests
```
yarn run test
```

### Lints and fixes files
```
yarn run lint
```

### Run your end-to-end tests
```
yarn run test:e2e
```

### Run your unit tests
```
yarn run test:unit
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
