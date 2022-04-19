<template>
  <div class="progressBar">
    <div class="progressBar__label">{{ this.title }} : {{ this.status }}</div>
    <progress class="progress is-primary" v-bind:value="this.value" max="100">
      {{ value }}%
    </progress>
  </div>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";

export default {
  name: "ProgressBar",
  data: function () {
    return {
      value: null,
      status: "Undefined",
      call: setInterval(this.updateStatus, 2000),
    };
  },
  props: {
    link: String,
    title: String,
  },
  methods: {
    updateStatus: function () {
      console.log(process.env.VUE_APP_API_ADDR + this.link);
      axios
        .get(process.env.VUE_APP_API_ADDR + this.link, {
          headers: {
            Authorization: cookies.get("Authorization"),
          },
        })
        .then((res) => {
          console.log(res);
        })
        .catch((err) => {
          console.log(err);
        });
    },
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
  .progress {
  }
}
</style>
