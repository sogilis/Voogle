<template>
  <div class="gallery">
    <h1 class="gallery__title">Gallery</h1>
    <PageNavigation
      :page="this.page"
      :is_last="is_last_page"
      :is_first="is_first_page"
      :attribute="this.attribute"
      :ascending="this.ascending"
      :status="this.status"
      :withSort="true"
      @pageChange="pageUpdate"
      @selectChange="selectUpdate"
    />
    <div v-if="this.status === 'Complete'">
      <button
        class="gallery__archive-button"
        :class="{ 'gallery__archive-button--cancel': this.enable_archive }"
        @click="
          this.enable_archive = !this.enable_archive;
          this.enable_unarchive = false;
          this.enable_deletion = false;
        "
      >
        <i
          class="gallery__archive-button-icon fa-solid fa-box-archive"
          v-if="!this.enable_archive"
        ></i>
        <i class="gallery__archive-button-icon fa-solid fa-ban" v-else></i>
      </button>
    </div>
    <div v-if="this.status === 'Archive'">
      <button
        class="gallery__unarchive-button"
        :class="{ 'gallery__unarchive-button--cancel': this.enable_unarchive }"
        @click="
          this.enable_archive = false;
          this.enable_unarchive = !this.enable_unarchive;
          this.enable_deletion = false;
        "
      >
        <i
          class="gallery__unarchive-button-icon fa-solid fa-boxes-packing"
          v-if="!this.enable_unarchive"
        ></i>
        <i class="gallery__unarchive-button-icon fa-solid fa-ban" v-else></i>
      </button>
      <button
        class="gallery__delete-button"
        :class="{ 'gallery__delete-button--cancel': this.enable_deletion }"
        @click="
          this.enable_archive = false;
          this.enable_unarchive = false;
          this.enable_deletion = !this.enable_deletion;
        "
      >
        <i
          class="gallery__delete-button-icon fa-solid fa-trash-can"
          v-if="!this.enable_deletion"
        ></i>
        <i class="gallery__delete-button-icon fa-solid fa-ban" v-else></i>
      </button>
    </div>
    <span>{{ this.error }}</span>
    <div class="gallery__wrapper">
      <div
        class="gallery__miniature-container"
        v-for="(video, index) in videos"
        :key="index"
      >
        <VideoMiniature
          :id="video.id"
          :title="video.title"
          :coverlink="video.coverlink"
          :enable_archive="this.enable_archive"
          :enable_unarchive="this.enable_unarchive"
          :enable_deletion="this.enable_deletion"
          @refreshResponse="this.refreshPage"
        ></VideoMiniature>
      </div>
    </div>
    <PageNavigation
      :page="this.page"
      :is_last="is_last_page"
      :is_first="is_first_page"
      :attribute="this.attribute"
      :ascending="this.ascending"
      :status="this.status"
      :withSort="false"
      @pageChange="pageUpdate"
      @selectChange="selectUpdate"
    />
  </div>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";
import VideoMiniature from "@/components/VideoMiniature";
import PageNavigation from "@/components/PageNavigation";

export default {
  name: "VideoGallery",
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
      status: "Complete",
      first_link: "",
      previous_link: "",
      next_link: "",
      last_link: "",
      enable_archive: false,
      enable_unarchive: false,
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
      return `api/v1/videos/list/${this.attribute}/${this.ascending}/${this.page}/${this.limit}/${this.status}`;
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
      this.status = payload.status;
      this.enable_archive = false;
      this.enable_unarchive = false;
      this.enable_deletion = false;
      this.page = 1;
      this.update(this.path);
    },
    refreshPage: function (payload) {
      if (!payload.error) {
        this.update(this.path);
      } else {
        this.error = payload.error;
        this.errored = true;
      }
    },
  },
  mounted() {
    this.update(this.path);
  },
  components: { VideoMiniature, PageNavigation },
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
  &__archive-button {
    opacity: 0.7;
    position: absolute;
    top: 20px;
    right: 20px;
    background-color: dimgray;
    color: white;
    border: 2px solid lightgray;
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
      background-color: red;
      opacity: 0.7;
      &:hover {
        opacity: 1;
      }
    }
  }
  &__unarchive-button {
    opacity: 0.7;
    position: absolute;
    top: 20px;
    right: 60px;
    background-color: green;
    color: white;
    border: 2px solid lightgray;
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
      background-color: red;
      opacity: 0.7;
      &:hover {
        opacity: 1;
      }
    }
  }
  &__delete-button {
    opacity: 0.7;
    position: absolute;
    top: 20px;
    right: 20px;
    background-color: red;
    color: white;
    border: 2px solid lightgray;
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
      background-color: red;
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
