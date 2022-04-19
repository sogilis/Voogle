<template>
  <div class="progressBar">
    <div class="progressBar__label">{{ this.title }} :</div>
    <progress class="progress is-primary" v-bind:value="this.value" max="100">
      {{ value }}%
    </progress>
  </div>
</template>

<script>
export default {
  name: "ProgressBar",
  data: function () {
    return {};
  },
  computed: {
    value: function () {
      switch (this.status) {
        case 0:
          return 10;
        case 200:
          return 100;
        default:
          return "";
      }
    },
  },
  props: {
    link: String,
    title: String,
    status: Number,
  },
  methods: {},
};
</script>

<style scoped lang="scss">
.progressBar {
  display: flex;
  width: 50%;
  flex-direction: column;
  justify-content: center;
  align-items: left;
  padding: 20px;
  .progress {
  }
}
</style>
