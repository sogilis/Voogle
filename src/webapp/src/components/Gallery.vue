<template>
  <div class="gallery">
    <h1 class="gallery__title">Gallery</h1>
    <PageNavigation
      :page="this.page"
      :is_last="is_last_page"
      :is_first="is_first_page"
      :attribute="this.attribute"
      :ascending="this.ascending"
      :withSort="true"
      @pageChange="pageUpdate"
      @selectChange="selectUpdate"
    />
    <button
      class="gallery__delete-button"
      :class="{ 'gallery__delete-button--cancel': this.enable_deletion }"
      @click="this.enable_deletion = !this.enable_deletion"
    >
      <i
        class="gallery__delete-button-icon fa-solid fa-trash-can"
        v-if="!this.enable_deletion"
      ></i>
      <i class="gallery__delete-button-icon fa-solid fa-ban" v-else></i>
    </button>
    <span>{{ this.error }}</span>
    <div class="gallery__wrapper">
      <div
        class="gallery__miniature-container"
        v-for="(video, index) in videos"
        :key="index"
      >
        <Miniature
          v-bind:id="video.id"
          v-bind:title="video.title"
          v-bind:enable_deletion="this.enable_deletion"
          @deletionResponse="this.handleDeletion"
        ></Miniature>
      </div>
    </div>
    <PageNavigation
      :page="this.page"
      :is_last="is_last_page"
      :is_first="is_first_page"
      :attribute="this.attribute"
      :ascending="this.ascending"
      :withSort="false"
      @pageChange="pageUpdate"
      @selectChange="selectUpdate"
    />
  </div>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";
import Miniature from "@/components/Miniature";
import PageNavigation from "@/components/PageNavigation";

export default {
  name: "Gallery",
  data: function () {
    return {
      videos: [],
      loading: true,
      errored: false,
      error: "",
      attribute: "upload_date",
      ascending: false,
      page: 1,
      last_page: 1,
      limit: 10,
      first_link: "",
      previous_link: "",
      next_link: "",
      last_link: "",
      enable_deletion: false,
    };
  },
  computed: {
    is_last_page: function () {
      return this.page == this.last_page;
    },
    is_first_page: function () {
      return this.page == 1;
    },
    path: function () {
      return `api/v1/videos/list/${this.attribute}/${this.ascending}/${this.page}/${this.limit}`;
    },
  },
  methods: {
    update: function (link) {
      axios
        .get(process.env.VUE_APP_API_ADDR + link, {
          headers: {
            Authorization: cookies.get("Authorization"),
          },
        })
        .then((response) => {
          this.videos = response.data["videos"];
          this.last_page = response.data["_lastpage"];
          this.first_link = response.data["_links"]["first"]["href"];
          this.previous_link = this.next_link = this.last_link = undefined;
          if (response.data["_links"]["previous"]) {
            this.previous_link = response.data["_links"]["previous"]["href"];
          }
          if (response.data["_links"]["next"]) {
            this.next_link = response.data["_links"]["next"]["href"];
            this.last_link = response.data["_links"]["last"]["href"];
          }
        })
        .catch((error) => {
          this.error = error;
          this.errored = true;
        });
    },
    pageUpdate: function (payload) {
      switch (payload.page) {
        case "first":
          this.update(this.first_link);
          this.page = 1;
          break;
        case "previous":
          this.update(this.previous_link);
          this.page = this.page - 1;
          break;
        case "next":
          this.update(this.next_link);
          this.page = this.page + 1;
          break;
        case "last":
          this.update(this.last_link);
          this.page = this.last_page;
          break;
      }
    },
    selectUpdate: function (payload) {
      this.attribute = payload.attribute;
      this.ascending = payload.ascending;
      this.page = 1;
      this.update(this.path);
    },
    handleDeletion: function (payload) {
      if (!payload.error) {
        this.update(this.path);
      } else {
        this.error = payload.error;
      }
    },
  },
  mounted() {
    this.update(this.path);
  },
  components: { Miniature, PageNavigation },
};
</script>

<style scoped lang="scss">
.gallery {
  position: relative;
  text-align: center;
  &__title {
    font-size: 1.5em;
    font-weight: bold;
    padding-top: 1em;
  }
  &__delete-button {
    opacity: 0.7;
    position: absolute;
    top: 20px;
    right: 20px;
    background-color: red;
    color: white;
    border: 1px solid black;
    font-size: 1.2rem;
    border-radius: 0.3em;
    &-icon {
      height: 1.2rem;
      width: 1.2rem;
    }
    &:hover {
      opacity: 1;
    }
    &--cancel {
      background-color: green;
      opacity: 0.7;
      &:hover {
        opacity: 1;
      }
    }
  }
  &__wrapper {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    grid-gap: 30px;
    padding: 1em;
  }
  &__miniature_container {
    height: 200px;
    width: 100%;
    max-width: 250px;
  }
}
</style>
