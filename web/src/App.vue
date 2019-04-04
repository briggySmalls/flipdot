<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <Login v-bind="client"/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Login from './components/Login.vue';
import { Client } from './ts/client';

const SERVER = 'https://jimsflipdot.hopto.org';

@Component({
  components: {
    Login,
  },
})
export default class App extends Vue {
  private client: Client;

  constructor() {
    super();
    this.client = new Client(SERVER);
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
