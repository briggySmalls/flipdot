<template>
  <form v-on:submit.prevent="sendMessage">
    <div class="form-group row">
        <label for="text" class="col-sm-2 col-form-label">Sender:</label>
        <b-form-input
          v-model="sender"
          placeholder="Sam"
          type="text"
          class="col-sm-10 form-control message-sender"
          :state="sender.length > 0"
          trim
          required>
        </b-form-input>
    </div>
    <div class="form-group row">
        <label for="text" class="col-sm-2 col-form-label">Message:</label>
        <b-form-textarea
          v-model="message"
          :state="message.length > 0"
          placeholder="Message here, try to keep it short!"
          rows="3"
          max-rows="4"
          class="col-sm-10 form-control message-text"
          trim
          required>
        </b-form-textarea>
    </div>
    <button type="submit" class="btn btn-primary message-submit">Send</button>
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
      if (this.client.error && this.client.error.code === grpc.Code.Unauthenticated) {
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
