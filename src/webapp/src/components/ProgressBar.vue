<template>
  <div class="progressBar">
    <div class="progressBar__label">
      {{ this.title }} : {{ this.statusArray[this.status - 1] }}
    </div>
    <progress class="progress is-primary" v-bind:value="this.value" max="100">
      {{ value }}%
    </progress>
  </div>
</template>

<script>
export default {
  name: "ProgressBar",
  data: function () {
    return {
      //List of available status.
      //Ensure status are provided in the correct order.
      statusArray: ["Uploading", "Uploaded", "Encoding", "Complete"],
    };
  },
  computed: {
    value: function () {
      return Math.round((100 / this.statusArray.length) * this.status);
    },
  },
  props: {
    status: Number,
    title: String,
  },
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
}
</style>
