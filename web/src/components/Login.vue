<template>
  <div>
    <b-alert class="prewrap" variant="danger" v-bind:show="client.error !== null">{{ client.error }}</b-alert>
    <b-form v-on:submit.prevent="authenticate">
      <b-form-group
        label="Password:"
        label-for="password-field">
          <b-form-input
            id="password-field"
            v-model="password"
            type="password"
            required>
          </b-form-input>
      </b-form-group>
      <b-button type="submit" variant="primary">Authorize</b-button>
    </b-form>
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
</style>
