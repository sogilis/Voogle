<template>
  <div>
    <h1>Gallery</h1>
    <div class="wrapper">
      <Video v-for="video in videos" :key="video" v-bind:video="video" />
    </div>
    <div class="button-generator">
      <button v-on:click="generateVideo">Generate video</button>
    </div>
  </div>
</template>

<script>
import Video from '@/components/Video.vue';

export default {
  name: 'Gallery',
  data: function () {
    return {
      videos: []
    }
  },
  methods: {
    generateVideo: function (event) {
      this.videos = { title: 'Video ' }
    }
  },
  components: {
    Video
  }
}
</script>

<style scoped lang="scss">
.wrapper {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  grid-gap: 10px;
  grid-auto-rows: minmax(100px, auto);
}

.button-generator {
  margin: 1%;
}
</style>
