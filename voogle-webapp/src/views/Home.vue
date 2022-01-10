<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <Gallery/>
  </div>
</template>

<script>
import Gallery from '@/components/Gallery.vue'

export default {
  name: 'Home',
  components: {
    Gallery
  }
}
</script>
