<template>
  <article @click="goToVideo" class="miniature">
    <figure class="minitature__preview">
      <img
        :src="'https://sogilis.com/wp-content/uploads/2021/09/logo_sogilis_alone.svg'"
        alt="video miniature"
      />
    </figure>
    <div class="miniature__title">{{ this.title }}</div>
    <button
      class="miniature__delete-button"
      @click.stop="this.delete()"
      v-if="enable_deletion"
    >
      <i class="fa-solid fa-trash-can"></i>
    </button>
  </article>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";

export default {
  name: "Miniature",
  props: {
    title: String,
    id: String,
    enable_deletion: Boolean,
  },
  methods: {
    goToVideo: function () {
      this.$router.push({ path: `/watch/${this.id}` });
    },
    delete: function () {
      axios
        .delete(
          process.env.VUE_APP_API_ADDR + `/api/v1/videos/${this.id}/delete`,
          {
            headers: {
              Authorization: cookies.get("Authorization"),
            },
          }
        )
        .then(() => {
          this.$emit("deletionResponse", {});
        })
        .catch((error) => {
          this.$emit("deletionResponse", { error: error });
        });
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
    transform: scale(1.05);

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

  &__delete-button {
    position: absolute;
    right: -1px;
    top: -1px;
    height: 24px;
    width: 24px;
    padding: 3px;
    background-color: red;
    color: white;
    border: none;
    border-radius: 5px;
  }
}

img {
  object-fit: cover;
}
</style>
