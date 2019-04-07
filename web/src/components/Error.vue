<template>
  <div id="message">{{ message }}</div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { grpc } from '@improbable-eng/grpc-web';

@Component
export default class Error extends Vue {
  @Prop() private error!: any;

  // Computed property for error message
  get message(): string {
    if (this.error === null) {
        return '';
    } else if (this.error.code === grpc.Code.Unauthenticated) {
        return `Authentication Error\n${this.error.message}\nPlease log in again`;
    }
    return this.error.message;
  }
}
</script>

<style>
#message {
    white-space: pre-wrap;
}
</style>