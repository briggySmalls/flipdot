<template>
  <Page v-bind:title="title" v-bind:text="text">
    <b-alert class="prewrap" v-bind:variant="alertVariant" show>{{ message }}</b-alert>
    <b-button v-on:click="newMessage" variant="primary" block>Start again</b-button>
  </Page>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Client } from '../ts/client';
import Page from './Page.vue';
import { grpc } from '@improbable-eng/grpc-web';

@Component({
  components: {
    Page,
  },
})
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

  // Page title
  get title(): string { return 'All done'; }

  // Page text
  get text(): string {
    return '';
  }
}
</script>

<style>
</style>
