<template>
  <form v-on:submit.prevent="sendMessage">
    <div class="form-group">
        <label for="text">Sender:</label>
        <input v-model="sender" type="text" class="form-control" required>
    </div>
    <div class="form-group">
        <label for="text">Message:</label>
        <input v-model="message" type="text" class="form-control" required>
    </div>
    <button type="submit" class="btn btn-primary">Send</button>
  </form>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Client } from '../ts/client';
import { grpc } from '@improbable-eng/grpc-web';

@Component
export default class Message extends Vue {
  // Name of sender
  public sender: string = '';
  // Message text
  public message: string = '';

  // Client used to communicate with flipdot display
  @Prop() private client!: Client;

  // Application state machine
  @Prop() private fsm!: any;

  public sendMessage() {
    // Send message to the server
    this.client.sendTextMessage(this.sender, this.message, (response) => {
      if (this.client.error === grpc.Code.Unauthenticated) {
        // Token has expired or something weirder: go back to login
        this.fsm.send('REAUTH');
        return;
      }
      // We have sent the message
      this.fsm.send('SENT');
    });
  }
}
</script>

<style>

</style>
