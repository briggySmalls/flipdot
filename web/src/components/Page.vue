<template>
  <div>
    <h2>{{ title }}</h2>
    <div class="mb-4">
      {{ text }}
    </div>
    <slot/>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Page extends Vue {
  // Title of the page
  @Prop() public title!: string;

  // Text of the page
  @Prop() public text!: string;
}
</script>

<style>
</style>
