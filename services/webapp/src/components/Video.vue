<template>
  <div>
    <h1>{{ video.title }}</h1>
    <video class="video-js vjs-theme-forest" :data-id="video.title" controls>
      <source :src="'https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8'" />
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
    const player = videojs(
      document.querySelector("video[data-id='" + this.video.title + "']")
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
