<template>
  <div>
    <b-alert class="prewrap" v-bind:variant="alertVariant" show>{{ message }}</b-alert>
    <b-button v-on:click="newMessage" variant="primary" block>Send a message</b-button>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Client } from '../ts/client';
import { grpc } from '@improbable-eng/grpc-web';

@Component
export default class Result extends Vue {
  // Client used to communicate with flipdot display
  @Prop() private client!: Client;

  // Application state machine
  @Prop() private fsm!: any;

  // Transition app back to message form
  public newMessage() {
    // Send 'new' event to state machine
    this.fsm.send('NEW');
  }

  get alertVariant(): string {
    // If there is no error, we have succeeded
    if (this.client.error === null) {
        return 'success';
    }
    // If there is an error, we have failed
    return 'danger';
  }

  // Computed property for error message
  get message(): string {
    // If there is no error, we display no message
    if (this.client.error === null) {
        return 'Message sent!';
    }
    // If there is an error, it must be to do with authentication
    return this.client.error;
  }
}
</script>

<style>
</style>
