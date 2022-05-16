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
      //List of available status.
      //Ensure status are provided in the correct order.
      statusArray: ["Uploading", "Uploaded", "Encoding"],
      requestLoop: setInterval(this.updateStatus, 500),
    };
  },
  props: {
    link: String,
    title: String,
  },
  beforeUnmount() {
    clearInterval(this.requestLoop);
  },
  methods: {
    updateStatus: function () {
      axios
        .get(process.env.VUE_APP_API_ADDR + this.link, {
          headers: {
            Authorization: cookies.get("Authorization"),
          },
        })
        .then((res) => {
          var valueGain = Math.floor(100 / this.statusArray.length);
          this.status = res.data["status"];
          if (this.status == "Complete") {
            this.value = 100;
            clearInterval(this.requestLoop);
          } else {
            this.value =
              this.statusArray.findIndex((s) => s == this.status) * valueGain;
          }
        })
        .catch((err) => {
          this.value = null;
          if (err.response) {
            this.status = err.response.status + " : " + err.response.data;
          } else {
            this.status = "Error, check your services health.";
          }
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
