<template>
  <div class="watchview">
    <h1 class="watchview__title">WATCHING</h1>
    <h2 class="watchview__video-title">{{ this.title }} - {{ this.date }}</h2>
    <Video :videoId="this.id" />
  </div>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";
import Video from "@/components/Video.vue";

export default {
  name: "Watch.vue",
  data: function () {
    return {
      id: this.$route.params.id,
      title: "",
      date: "",
    };
  },
  mounted() {
    axios
      .get(process.env.VUE_APP_API_ADDR + `api/v1/videos/${this.id}/info`, {
        headers: {
          Authorization: cookies.get("Authorization"),
        },
      })
      .then((response) => {
        this.title = response.data["title"];
        this.date = response.data["uploaddate"];
      })
      .catch((error) => {
        this.title = error;
      });
  },
  components: {
    Video,
  },
};
</script>

<style scoped lang="scss">
.watchview {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  row-gap: 20px;

  &__title {
    font-size: 1.5em;
    font-weight: bold;
    padding-top: 1em;
  }

  &__video-title {
    font-size: 1em;
    font-weight: bold;
  }
}
</style>
