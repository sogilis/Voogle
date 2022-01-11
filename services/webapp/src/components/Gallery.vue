<template>
  <div>
    <h1>Gallery component</h1>
    <div class="wrapper">
      <Video v-for="video in videos" :key="video" v-bind:video="video" />
    </div>
    <h2>Test API</h2>
    <div v-for="currency in info" :key="currency">
      {{ currency.description }}:
      <span class="lighten">
        <span v-html="currency.symbol"></span>{{ currency.rate_float }}
      </span>
    </div>
    <div>
      <label for="video-generator">Number video generator:</label>
      <input v-model="nbVideo">
      <button v-on:click="generateVideo">Generate {{nbVideo}} video</button>
    </div>
    <div>
      <label for="increment">Increment:</label>
      <input v-model="increment">
      <h2 v-on:click="incrementClickCounter">Click counter {{clickCounter}}</h2>
    </div>
  </div>
</template>

<script>
import Video from '@/components/Video.vue'
import axios from 'axios'

export default {
  name: 'Gallery',
  props: {
    msg: String
  },
  data: function () {
    return {
      videos: [],
      nbVideo: 0,
      increment: 0,
      info: null,
      loading: true,
      errored: false
    }
  },
  mounted () {
    axios
      .get('https://api.coindesk.com/v1/bpi/currentprice.json')
      .then(response => {
        (this.info = response.data.bpi)
      })
      .catch(error => {
        console.log(error)
        this.errored = true
      })
      .finally(() => {
        this.loading = false
      })
  },
  computed: {
    clickCounter () {
      return this.$store.state.clickCounter
    }
  },
  methods: {
    incrementClickCounter: function (event) {
      this.$store.commit('setClickCounter', this.clickCounter + Number(this.increment))
    },
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
