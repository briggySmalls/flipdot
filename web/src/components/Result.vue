<template>
  <div id="result">
    <div id="message">{{ message }}</div>
    <button v-on:click="newMessage">Send a message</button>
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
#message {
    white-space: pre-wrap;
}
</style>