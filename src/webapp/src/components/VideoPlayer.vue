""""""""""""""""""""
<template>
  <div>
    <video
      class="video-js vjs-theme-forest"
      :data-id="videoId"
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
import cookies from "js-cookie";

export default {
  props: {
    videoId: String,
    reqParameter: String,
  },
  mounted() {
    videojs.Vhs.xhr.beforeRequest = (options) => {
      options.headers = options.headers || {};
      options.headers.Authorization = cookies.get("Authorization");
      options.uri += this.reqParameter;
      return options;
    };
    const player = videojs(
      document.querySelector("video[data-id='" + this.videoId + "']")
    );
    player.src(
      process.env.VUE_APP_API_ADDR +
        "api/v1/videos/" +
        this.videoId +
        "/streams/master.m3u8"
    );
    player.hlsQualitySelector({
      displayCurrentQuality: true,
    });
    this.videoplayer = player;
  },
};
</script>

<style scoped lang="scss">
.video-js {
  display: block;
  margin: 0 auto;
}
</style>
