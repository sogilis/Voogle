<template>
  <div>
    <h1>Gallery</h1>
    <div class="wrapper">
      <Video v-for="video in videos" :key="video" v-bind:video="video" />
    </div>
  </div>
</template>

<script>
import axios from "axios";

import Video from "@/components/Video.vue";

export default {
  name: "Gallery",
  data: function () {
    return {
      videos: [],
      loading: true,
      errored: false,
      error: "",
    };
  },
  mounted() {
    axios
      .get(process.env.VUE_APP_API_ADDR + "/api/v1/videos/list", {
        auth: {
          username: process.env.VUE_APP_API_USER,
          password: process.env.VUE_APP_API_PWD,
        },
      })
      .then((response) => {
        this.videos = response.data.data;
      })
      .catch((error) => {
        this.error = error;
        this.errored = true;
      })
      .finally(() => {
        this.loading = false;
      });
  },
  components: {
    Video,
  },
};
</script>

<style scoped lang="scss">
.wrapper {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  grid-gap: 10px;
  grid-auto-rows: minmax(100px, auto);
}
</style>
