<template>
  <div>
    <h1>{{ video.title }}</h1>
    <video
      class="video-js vjs-theme-forest"
      :data-id="video.id"
      controls
    ></video>
  </div>
</template>

<script>
import "video.js/dist/video-js.css";
import "@videojs/themes/dist/forest/index.css";
import videojs from "video.js";
import "videojs-hls-quality-selector";
import "videojs-contrib-quality-levels";

export default {
  props: {
    video: Object,
  },
  mounted() {
    videojs.Hls.xhr.beforeRequest = function (options) {
      options.headers = options.headers || {};
      options.headers.Authorization =
        "Basic " +
        btoa(process.env.VUE_APP_API_USER + ":" + process.env.VUE_APP_API_PWD);
      return options;
    };
    const player = videojs(
      document.querySelector("video[data-id='" + this.video.id + "']")
    );
    player.src(
      process.env.VUE_APP_API_ADDR +
        "/api/v1/videos/" +
        this.video.id +
        "/streams/master.m3u8"
    );
    player.hlsQualitySelector({
      displayCurrentQuality: true,
    });
  },
};
</script>

<style scoped lang="scss">
.video-js {
  display: block;
  margin: 0 auto;
}
</style>
