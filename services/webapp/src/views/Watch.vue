<template>
  <h1>WATCHING</h1>
  <Video
    :key="{ id: this.$route.params.id, title: this.$route.params.id }"
    v-bind:video="{ id: this.$route.params.id, title: this.$route.params.id }"
  />
</template>

<script>
import Video from "@/components/Video.vue";

export default {
  name: "Watch.vue",
  components: {
    Video,
  },
};
</script>

<style scoped></style>
