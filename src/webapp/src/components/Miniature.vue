<template>
  <article @click="goToVideo" class="miniature">
    <figure class="minitature__preview">
      <img
        :src="'https://sogilis.com/wp-content/uploads/2021/09/logo_sogilis_alone.svg'"
        alt="video miniature"
      />
    </figure>
    <div class="miniature__title">{{ this.title }}</div>
  </article>
</template>

<script>
export default {
  name: "Miniature",
  props: {
    title: String,
    id: String,
  },
  methods: {
    goToVideo() {
      this.$router.push({ path: `/watch/${this.id}` });
    },
  },
};
</script>

<style scoped lang="scss">
.miniature {
  $block-element: &;
  width: 100%;
  max-width: 250px;
  height: 200px;
  border: 1px solid black;
  border-radius: 5px;
  overflow: hidden;
  transition: all 400ms;
  position: relative;

  &:hover {
    cursor: pointer;
    transform: scale(1.1);

    #{$block-element}__title {
      max-height: 5em;
    }
  }

  &__preview {
    padding: 5px;
    height: 100%;
    width: 100%;
  }

  &__title {
    position: absolute;
    padding: 0px;
    bottom: 0px;
    font-size: 1.3em;
    max-height: 1.4em;
    width: 100%;
    text-align: center;
    border-radius: 5px;
    background-color: #e9e9e9;
    transition: max-height 400ms;
  }
  figure {
    padding: 5px;
    height: 100%;
    width: 100%;
  }
  &__title {
    position: absolute;
    padding: 0px;
    bottom: 0px;
    font-size: 1.3em;
    max-height: 1.4em;
    width: 100%;
    text-align: center;
    border-radius: 5px;
    background-color: #e9e9e9;
    transition: max-height 400ms;
  }
}

img {
  object-fit: cover;
}
</style>
