<template>
  <div>
    <h1>Gallery component</h1>
    <div class="wrapper">
      <Video v-for="video in videos" :key="video" v-bind:video="video" />
    </div>
    <div>
      <label for="video-generator">Number video generator:</label>
      <input v-model="nbVideo">
      <button v-on:click="generateVideo">Generate {{nbVideo}} video</button>
    </div>
  </div>
</template>

<script>
import Video from '@/components/Video.vue'

export default {
  name: 'Gallery',
  props: {
    msg: String
  },
  data: function () {
    return {
      videos: [],
      nbVideo: 0,
      info: null,
      loading: true,
      errored: false
    }
  },
  computed: {
  },
  methods: {
    generateVideo: function (event) {
      for (let i = 0; i < this.nbVideo; i++) {
        this.videos.push({ title: 'Video ' + i, description: 'Video Description ' + i })
      }
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
  grid-template-columns: repeat(4, 1fr);
  grid-gap: 10px;
  grid-auto-rows: minmax(100px, auto);
}
</style>
