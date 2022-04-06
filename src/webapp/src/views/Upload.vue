<template>
  <h2>UPLOAD</h2>
  <form
    class="form_upload flex_center"
    @submit.prevent="submitFile()"
    v-if="this.status === 'None'"
  >
    <div id="video_box" :class="{ active: fileSelected }">
      <input
        id="file"
        type="file"
        ref="file"
        @change="handleFileUpload()"
        accept="video/*"
      />
      <div class="in_video_box flex_center" v-if="!fileSelected">
        <!-- Div handling drag-events when supported -->
        <div
          id="drag_box"
          v-if="dragEnabled"
          @drop.prevent.stop="handleDrop"
          @dragleave.prevent.stop="dragOnBox(false)"
          @dragend.prevent.stop="dragOnBox(false)"
          @dragenter.prevent.stop="dragOnBox(true)"
          @dragover.prevent.stop="dragOnBox(true)"
          @drag.prevent.stop=""
          @dragstart.prevent.stop=""
        ></div>
        <label
          id="box_text"
          for="file"
          @dragenter.prevent.stop="dragOnBox(true)"
        >
          <strong>Choose a file</strong>
          <span v-if="isAdvancedUpload()"
            ><br />
            or drag it here</span
          >.
        </label>
      </div>
      <div class="in_video_box flex_center" v-else>
        <img
          id="video_preview"
          src="../assets/uxwing-camera.png"
          alt="Video Set"
        />
        <div class="title_box flex_center">
          <input
            type="text"
            placeholder="Enter a Title"
            v-model="title"
            required
          />
        </div>
      </div>
    </div>
    <span class="flex_center">
      <button type="submit" class="button is-primary" :disabled="!fileSelected">
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
  <div v-else>
    <h5>{{ this.status }}</h5>
    <span>Uploading and transforming a video can be a lengthy process</span>
    <br />
    <span @click="retry()">Retry</span>
  </div>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";

export default {
  name: "Upload",
  data: function () {
    return {
      status: "None",
      title: "",
      file: "",
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
      console.log("Submitting");
      if (this.title == "") {
        alert("Please enter a title.");
        return;
      }
      // Simple way to keep the user aware about what is happening
      this.status = "Uploading";

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
          if (res.status === 200) {
            this.status = "Uploaded";
          } else {
            this.status = "Failed - " + res.statusText;
          }
        })
        .catch((err) => {
          this.status = err;
        });
    },
    handleFileUpload() {
      // v-model doesn't support file form type
      this.file = this.$refs.file.files[0];
    },
    retry: function () {
      this.status = "None";
      this.title = "";
      this.file = "";
    },
    isAdvancedUpload: function () {
      var div = document.createElement("div");
      return (
        ("draggable" in div || ("ondragstart" in div && "ondrop" in div)) &&
        "FormData" in window &&
        "FileReader" in window
      );
    },
    handleDrop: function (e) {
      this.file = e.dataTransfer.files[0];
      this.dragOnBox(false);
    },
    dragOnBox: function (bool) {
      var back = document.getElementById("video_box");
      var front = document.getElementById("drag_box");
      if (bool) {
        back.classList.add("is_dragged_over");
        front.classList.add("is_dragged_over");
      } else {
        back.classList.remove("is_dragged_over");
        front.classList.remove("is_dragged_over");
      }
    },
  },
};
</script>

<style scoped lang="scss">
.form_upload {
  margin: auto;
  width: fit-content;
  flex-direction: column;
  row-gap: 10px;
}

#video_box {
  float: right;
  height: 200px;
  width: 200px;
  border: 5px rgb(175, 175, 255) dashed;
  border-radius: 10px;
  overflow: hidden;
  transition: background-color 400ms;
  &.is_dragged_over {
    border-style: solid;
    background-color: rgb(218, 218, 255);
  }
  &.active {
    border: 2px rgb(175, 175, 255) solid;
  }
}

.in_video_box {
  height: inherit;
  width: inherit;
  flex-direction: column;
}

#drag_box {
  height: inherit;
  width: inherit;
  position: absolute;
  &.is_dragged_over {
    z-index: 2;
  }
}

.label {
  font-size: larger;
  font-weight: bold;
}

#box_text {
  z-index: 1;
  strong:hover {
    text-decoration: underline;
  }
}

#file {
  display: none;
}

#video_preview {
  width: 60%;
  height: 60%;
  margin: 10% 20%;
}

#cancel {
  margin-left: 20px;
  background-color: red;
}

.title_box {
  height: 20%;
  width: inherit;
  background-color: black;
  position: bottom;
}

#show_title {
  flex-grow: 1;
  color: white;
  font-size: large;
  width: 80%;
}

input {
  border: none;
  height: 100%;
  width: 80%;
  background: inherit;
  color: white;
  font-size: large;
  outline: none;
}

button {
  width: 90px;
  margin: 0 10px;
  i {
    margin-left: 5px;
  }
}
</style>
