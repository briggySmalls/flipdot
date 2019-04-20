<template>
  <b-form v-on:submit.prevent="sendMessage">
    <b-form-group
      label="Sender:"
      label-for="sender-field">
      <b-form-input
        id="sender-field"
        v-model="sender"
        placeholder="Your name"
        type="text"
        :state="sender.length > 0"
        trim
        required>
      </b-form-input>
    </b-form-group>
    <b-form-group
      label="Message:"
      label-for="text-field">
      <b-form-textarea
        id="text-field"
        v-model="message"
        :state="message.length > 0"
        placeholder="Your message here, try to keep it short!"
        rows="3"
        max-rows="4"
        trim
        required>
      </b-form-textarea>
    </b-form-group>
    <b-button id="message-submit" type="submit" variant="primary" block>Send</b-button>
  </b-form>
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
      // Clear the form
      this.reset();
      if (this.client.error && this.client.error.code === grpc.Code.Unauthenticated) {
        // Token has expired or something weirder: go back to login
        this.fsm.send('REAUTH');
        return;
      }
      // We have sent the message
      this.fsm.send('SENT');
    });
  }

  private reset() {
    this.sender = '';
    this.message = '';
  }
}
</script>

<style>
</style>
