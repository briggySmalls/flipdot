<template>
  <div>
    <div v-show="client.error !== null" class="prewrap alert alert-danger">{{ message }}</div>
    <form v-on:submit.prevent="authenticate">
      <div class="form-group row">
          <label for="text" class="col-sm-2 col-form-label">Password:</label>
          <input v-model="password" type="password" class="col-sm-10 form-control login-password" required>
      </div>
      <button type="submit" class="btn btn-primary login-submit">Authorize</button>
    </form>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Client } from '../ts/client';
import { grpc } from '@improbable-eng/grpc-web';
import {AuthenticateRequest} from '../generated/flipapps_pb';

@Component
export default class Login extends Vue {
  // Password bound to the view
  public password: string = '';

  // Client used to communicate with flipdot display
  @Prop() private client!: Client;

  // Application state machine
  @Prop() private fsm!: any;

  // Attempt to authenticate using client
  public authenticate() {
    // Authenticate with the client
    this.client.authenticate(this.password, (response) => {
      if (this.client.error === null) {
        // We authenticated correctly, transition to sending a message
        this.fsm.send('AUTH');
      }
    });
  }

  // Computed property for error message
  get message(): string {
    // If there is no error, we display no message
    if (this.client.error === null) {
        return '';
    }
    // If there is an error, it must be to do with authentication
    return `Authentication Error\n${this.client.error.message}\nPlease log in again`;
  }
}
</script>

<style>
@import "../assets/common.scss";
</style>
