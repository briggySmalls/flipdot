<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <Error v-show="error != null" v-bind:error="error"/>
    <Login v-show="!isAuthenticated" v-bind:client="client"/>
    <Message v-show="isAuthenticated" v-bind:client="client"/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Login from './components/Login.vue';
import Message from './components/Message.vue';
import Error from './components/Error.vue';
import { Client } from './ts/client';

const SERVER = 'https://jimsflipdot.hopto.org';

@Component({
  components: {
    Login,
    Message,
    Error,
  },
})
export default class App extends Vue {
  // Client to be shared between components
  public client: Client;

  constructor() {
    // Call super
    super();
    // Create a client
    this.client = new Client(SERVER);
  }

  get isAuthenticated() {
    return this.client.isAuthenticated;
  }

  get error() {
    return this.client.error;
  }
}
</script>

<style lang="scss">
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
