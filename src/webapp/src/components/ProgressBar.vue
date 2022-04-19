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
      call: setInterval(this.updateStatus, 500),
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
          this.status = res.data["status"]
          switch (this.status){
            case ("VIDEO_STATUS_UPLOADING"): {
              this.value = 0;
              break;
            }
            case ("VIDEO_STATUS_UPLOADED"): {
              this.value = 25;
              break;
            }
            case ("VIDEO_STATUS_ENCODING"): {
              this.value = 60;
              break;
            }
            case ("VIDEO_STATUS_COMPLETE"): {
              this.value = 100;
              clearInterval(this.call);
              break;
            }
            default: {
              this.value = null;
            }
          }
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
