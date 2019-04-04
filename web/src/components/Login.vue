<template>
    <form v-on:submit.prevent="authenticate" id="auth-form">
      <div class="form-group">
          <label for="text">Password:</label>
          <input v-model="password" type="password" name="password" id="password" class="form-control" required>
      </div>
      <button type="submit" class="btn btn-primary">Authorize</button>
    </form>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Client } from '../ts/client';

@Component
export default class Login extends Vue {
  // Password bound to the view
  public password: string = '';

  @Prop() private client!: Client;

  public authenticate() {
    this.client.authenticate(this.password);
  }
}
</script>

<style>

</style>