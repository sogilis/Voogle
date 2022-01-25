<template>
  <div v-on:click="goto" class="miniature">
    <h3> {{ this.title }} </h3>
    <img :src="'https://sogilis.com/wp-content/uploads/2021/09/logo_sogilis_alone.svg'"
         alt="video miniature">
  </div>
</template>
<script>
export default {
  name: 'Miniature',
  props: {
    title: String,
    id: String,
  },
  methods: {
    goto() {
      this.$router.push({path: `/watch/${this.id}`});
    },
  },
};
</script>

<style scoped lang="scss">
.miniature {
  padding: 5px;
  width: 80%;
  display: inline-block;
  border: 1px solid black;
}

img {
  max-width: 250px;
  max-height: 150px;
}
</style>
