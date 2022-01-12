<template>
  <div>
    <h1>{{ video.title }}</h1>
    <video id="video-player" class="video-js vjs-theme-forest" controls>
      <source
        :src="'https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8'"
        type="application/x-mpegURL"
      />
    </video>
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
    const player = videojs('video-player');
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
