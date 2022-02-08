<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png" />
    <img alt="Sqsc" src="../assets/sqsc.jpeg" />
    <img alt="Sqsc" src="../assets/sqsc.jpeg" />
    <img alt="g" src="../assets/g.png" />
    <img alt="l" src="../assets/l.png" />
    <img alt="e" src="../assets/e.png" />
    <Gallery />
  </div>
</template>

<script>
import Gallery from "@/components/Gallery.vue";

export default {
  name: "Home",
  components: {
    Gallery,
  },
};
</script>

<style scoped lang="scss">
.home {
  margin-top: 2%;
}

img {
  width: 50px;
  height: 50px;
}
</style>
