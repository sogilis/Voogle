<template>
  <article @click="goToVideo" class="miniature">
    <figure class="minitature__preview">
      <img v-bind:src="this.coverSrc" alt="video miniature" />
    </figure>
    <div class="miniature__title">{{ this.title }}</div>
    <button
      class="miniature__archive-button"
      @click.stop="this.archive()"
      v-if="enable_archive"
    >
      <i class="fa-solid fa-box-archive"></i>
    </button>
    <button
      class="miniature__unarchive-button"
      @click.stop="this.unarchive()"
      v-if="enable_unarchive"
    >
      <i class="fa-solid fa-boxes-packing"></i>
    </button>
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
  name: "VideoMiniature",
  props: {
    title: String,
    id: String,
    coverlink: Object,
    enable_archive: Boolean,
    enable_unarchive: Boolean,
    enable_deletion: Boolean,
  },
  data: function () {
    return {
      coverSrc: undefined,
    };
  },
  mounted() {
    this.getCover();
  },
  methods: {
    goToVideo: function () {
      this.$router.push({ path: `/watch/${this.id}` });
    },
    getCover: function () {
      axios
        .get(process.env.VUE_APP_API_ADDR + this.coverlink["href"], {
          headers: {
            Authorization: cookies.get("Authorization"),
          },
        })
        .then((response) => {
          if (response.data.length > 0) {
            this.coverSrc = "data:image/jpeg;base64," + response.data;
          } else {
            this.coverSrc =
              "https://sogilis.com/wp-content/uploads/2021/09/logo_sogilis_alone.svg";
          }
        })
        .catch(() => {
          this.coverSrc =
            "https://sogilis.com/wp-content/uploads/2021/09/logo_sogilis_alone.svg";
        });
    },
    archive: function () {
      axios
        .put(
          process.env.VUE_APP_API_ADDR + `api/v1/videos/${this.id}/archive`,
          undefined,
          {
            headers: {
              Authorization: cookies.get("Authorization"),
            },
          }
        )
        .then(() => {
          this.$emit("refreshResponse", {});
        })
        .catch((error) => {
          this.$emit("refreshResponse", { error: error });
        });
    },
    unarchive: function () {
      axios
        .put(
          process.env.VUE_APP_API_ADDR + `api/v1/videos/${this.id}/unarchive`,
          undefined,
          {
            headers: {
              Authorization: cookies.get("Authorization"),
            },
          }
        )
        .then(() => {
          this.$emit("refreshResponse", {});
        })
        .catch((error) => {
          this.$emit("refreshResponse", { error: error });
        });
    },
    delete: function () {
      axios
        .delete(
          process.env.VUE_APP_API_ADDR + `api/v1/videos/${this.id}/delete`,
          {
            headers: {
              Authorization: cookies.get("Authorization"),
            },
          }
        )
        .then(() => {
          this.$emit("refreshResponse", {});
        })
        .catch((error) => {
          this.$emit("refreshResponse", { error: error });
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

  &__archive-button {
    position: absolute;
    right: -1px;
    top: -1px;
    height: 24px;
    width: 24px;
    padding: 3px;
    background-color: dimgray;
    color: white;
    border: none;
    border-radius: 5px;
  }
  &__unarchive-button {
    position: absolute;
    right: -1px;
    top: -1px;
    height: 24px;
    width: 24px;
    padding: 3px;
    background-color: green;
    color: white;
    border: none;
    border-radius: 5px;
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
