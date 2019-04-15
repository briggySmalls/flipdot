<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <div class="container">
      <Login id="login" v-show="state == 'login'" v-bind:client="client"  v-bind:fsm="fsm"/>
      <Message id="message" v-show="state == 'message'" v-bind:client="client" v-bind:fsm="fsm"/>
      <Result id="result" v-show="state == 'result'" v-bind:client="client" v-bind:fsm="fsm"/>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Login from './components/Login.vue';
import Message from './components/Message.vue';
import Result from './components/Result.vue';
import { Client } from './ts/client';
import { createFsm } from './ts/fsm';
import BootstrapVue from 'bootstrap-vue';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

// Use Bootstrap
Vue.use(BootstrapVue);

@Component({
  components: {
    Login,
    Message,
    Result,
  },
})
export default class App extends Vue {
  // Client to be shared between components
  public client: Client;
  // State machine for controlling state
  private fsm: any;

  constructor() {
    // Call super
    super();
    // Create a state machine
    this.fsm = createFsm()
      .onTransition((state) => console.log(state.value))
      .start();
    // Create a client
    this.client = new Client(process.env.VUE_APP_GRPC_SERVER_URL);
  }

  get state() {
    return this.fsm.state.value;
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
