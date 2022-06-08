<template>
  <div class="video__container">
    <video class="video-js vjs-theme-forest" ref="videoId" controls></video>
  </div>
</template>

<script>
import "video.js/dist/video-js.css";
import "@videojs/themes/dist/forest/index.css";
import videojs from "video.js";
import videojsqualityselector from "videojs-hls-quality-selector";
import "videojs-contrib-quality-levels";
import cookies from "js-cookie";

export default {
  props: {
    videoId: String,
    filterlist: String,
  },
  data() {
    return {
      videoPlayer: null,
      timestamp: 0,
    };
  },
  computed: {
    filteruri: function () {
      return this.filterlist;
    },
  },
  mounted() {
    this.initializePlayer();
    this.videoPlayer.responsive(false);
    this.videoPlayer.hlsQualitySelector = videojsqualityselector;
    this.videoPlayer.hlsQualitySelector({
      displayCurrentQuality: true,
    });
  },
  beforeUpdate() {
    this.videoPlayer.pause();
    this.timestamp = this.videoPlayer.currentTime();
    this.videoPlayer.reset();
    this.initializePlayer();
    this.videoPlayer.currentTime(this.timestamp);
    this.videoPlayer.play();
  },
  methods: {
    initializePlayer: function () {
      videojs.Hls.xhr.beforeRequest = (options) => {
        options.headers = options.headers || {};
        options.headers.Authorization = cookies.get("Authorization");
        options.uri += this.filteruri;
        return options;
      };
      var player = videojs(this.$refs.videoId);
      player.src(
        process.env.VUE_APP_API_ADDR +
          "api/v1/videos/" +
          this.videoId +
          "/streams/master.m3u8"
      );
      this.videoPlayer = player;
    },
  },
};
</script>

<style scoped lang="scss">
.video-js {
  display: block;
  margin: 0 auto;
  height: 480px;
  width: 853px;
}
</style>
