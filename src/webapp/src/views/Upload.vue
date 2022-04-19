<template>
  <div class="uploadpage">
    <h2 class="uploadpage__title">UPLOAD</h2>
    <form class="uploadpage__form" @submit.prevent="submitFile()">
      <UploadBox :title="this.file.name" @sendFile="handleFile" />
      <label class="uploadpage__form-label" for="videotitle"
        >Give your video a title : </label
      ><input
        class="uploadpage__form-input"
        id="videotitle"
        type="text"
        placeholder="Enter a Title"
        v-model="title"
        required
      />
      <span class="uploadpage__form-buttoncontainer">
        <button
          type="submit"
          class="button is-primary"
          :disabled="!fileSelected"
        >
          <span>Upload</span>
          <span><i class="fa-solid fa-upload"></i></span>
        </button>
        <button
          class="button is-danger is-outlined"
          :disabled="!fileSelected"
          @click.stop.prevent="retry()"
        >
          <span>Cancel</span>
          <span class="icon is-small"> <i class="fa-solid fa-xmark"></i></span>
        </button>
      </span>
    </form>
    <div v-for="(upload, index) in progressArray" :key="index">
      <ProgressBar :title="upload.title" :link="upload.link" />
    </div>
  </div>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";
import ProgressBar from "@/components/ProgressBar";
import UploadBox from "@/components/UploadBox";

export default {
  name: "Upload",
  components: {
    ProgressBar,
    UploadBox,
  },
  data: function () {
    return {
      title: "",
      file: "",
      progressArray: [],
    };
  },
  computed: {
    fileSelected: function () {
      return !this.file == "";
    },
    dragEnabled: function () {
      return this.isAdvancedUpload() && !this.fileSelected;
    },
  },
  methods: {
    submitFile: function () {
      // Creating a FormData to POST it as multipart FormData
      const formData = new FormData();
      formData.append("title", this.title);
      formData.append("video", this.file);

      axios
        .post(process.env.VUE_APP_API_ADDR + "api/v1/videos/upload", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: cookies.get("Authorization"),
          },
        })
        .then((res) => {
          // Creating a new progress bar showing video status
          this.progressArray.push({
            title: this.title,
            link: res.data["links"].find(element => element["rel"]=="Status")["href"],
          });
          this.retry();
        })
        .catch((err) => {
          console.log(err);
        });
    },
    retry: function () {
      this.file = "";
    },
    handleFile: function (payload) {
      this.file = payload.file;
    },
  },
};
</script>

<style scoped lang="scss">
.uploadpage {
  &__title {
    font-size: 1.5em;
    font-weight: bold;
    padding-top: 1em;
    text-align: center;
  }
  &__form {
    padding-top: 1em;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    row-gap: 1em;
    &-label {
      font-size: 1.1em;
    }
    &-input {
      padding: 5px 15px;
    }
  }
}

button {
  margin: 0px 10px;
}
</style>
