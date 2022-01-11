<template>
  <div>
    <h1>API test</h1>
      <h2>Test API</h2>
      <div v-for="currency in info" :key="currency">
        {{ currency.description }}:
        <span class="lighten">
          <span v-html="currency.symbol"></span>{{ currency.rate_float }}
        </span>
      </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  props: {
    video: Object
  },
  data: function () {
    return {
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
  }
}
</script>
