<template>
  <div>
    <h1>Gallery</h1>
    <div class="wrapper">
      <div
        class="miniature_container"
        v-for="(video, index) in videos"
        :key="index"
      >
        <Miniature v-bind:id="video.id" v-bind:title="video.title"></Miniature>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

import cookies from "js-cookie";
import Miniature from "@/components/Miniature";

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
      .get(process.env.VUE_APP_API_ADDR + "api/v1/videos/list", {
        headers: {
          Authorization: cookies.get("Authorization"),
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
  components: { Miniature },
};
</script>

<style scoped lang="scss">
.wrapper {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  grid-gap: 10px;
  grid-auto-rows: minmax(100px, auto);
}
</style>
