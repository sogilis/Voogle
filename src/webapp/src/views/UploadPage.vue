<template>
  <div class="uploadpage">
    <h2 class="uploadpage__title">UPLOAD</h2>
    <form class="uploadpage__form" @submit.prevent="submitFile()">
      <label class="uploadpage__form-label" for="videofile"
        >Choose your video file :
      </label>
      <UploadBox
        :title="this.file.name"
        :accepting="'video/*'"
        :refto="'video_file'"
        @sendFile="handleFile"
      />
      <label class="uploadpage__form-label" for="videocover"
        >Choose your video cover :
      </label>
      <UploadBox
        :title="this.cover.name"
        :accepting="'image/png'"
        :refto="'cover_file'"
        @sendFile="handleCover"
      />
      <label class="uploadpage__form-label" for="videotitle"
        >Give your video a title :
      </label>
      <input
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
      <div v-if="msg">{{ msg }}</div>
    </form>
    <div v-for="(upload, index) in progressArray" :key="index">
      <ProgressBar :title="upload.title" :status="upload.status" />
    </div>
  </div>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";
import ProgressBar from "@/components/ProgressBar";
import UploadBox from "@/components/UploadBox";

export default {
  name: "UploadPage",
  components: {
    ProgressBar,
    UploadBox,
  },
  data: function () {
    return {
      title: "",
      file: "",
      cover: "",
      progressArray: [],
      msg: "",
      ws: "",
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
  mounted: function () {
    if (this.ws == "") {
      try {
        console.log(`Connecting to ws://${window.location.host}/ws`);
        this.ws = new WebSocket(`ws://${window.location.host}/ws`);

        this.ws.onopen = () => {};

        this.ws.onmessage = (event) => {
          try {
            let data = JSON.parse(event.data);
            let index = this.progressArray.findIndex(
              (upload) => upload["title"] == data["title"]
            );
            this.progressArray[index]["status"] = data["status"];
          } catch {
            this.msg = event.data;
          }
        };
      } catch (error) {
        this.msg = error;
      }
    }
  },
  beforeUnmount: function () {
    try {
      this.ws.close(1000, "Leaving upload page.");
    } catch (error) {
      this.msg = error;
    }
  },
  methods: {
    submitFile: function () {
      // Creating a FormData to POST it as multipart FormData
      this.progressArray.push({
        title: this.title,
        status: "Undefined",
      });
      const formData = new FormData();
      formData.append("title", this.title);
      formData.append("video", this.file);
      formData.append("cover", this.cover);
      try {
        this.ws.send(this.title);
      } catch {
        this.msg = "Could not subscribe to updates for '" + this.title + "'.";
      }

      axios
        .post(process.env.VUE_APP_API_ADDR + "api/v1/videos/upload", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: cookies.get("Authorization"),
          },
        })
        .then(() => {
          this.retry();
        })
        .catch((err) => {
          this.msg = err;
        });
    },
    retry: function () {
      this.title = "";
      this.file = "";
      this.cover = "";
    },
    handleFile: function (payload) {
      this.file = payload.file;
    },
    handleCover: function (payload) {
      this.cover = payload.file;
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
