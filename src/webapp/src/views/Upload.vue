<template>
  <h2 v-on:drop.prevent.stop>UPLOAD</h2>
  <form class="form_upload" v-if="this.status === 'None'">
    <div id="video_box">
      <input
        id="file"
        type="file"
        ref="file"
        @change="handleFileUpload()"
        accept="video/*"
        required
      />
      <!-- Div handling drag-events when supported -->
      <div
        id="drag_box"
        v-if="isAdvancedUpload()"
        @drop.prevent.stop="handleDrop"
        @dragleave.prevent.stop="dragOnBox(false)"
        @dragend.prevent.stop="dragOnBox(false)"
        @dragenter.prevent.stop="dragOnBox(true)"
        @dragover.prevent.stop="dragOnBox(true)"
        @drag.prevent.stop=""
        @dragstart.prevent.stop=""
      ></div>
      <label id="box_text" for="file">
        <strong>Choose a file</strong>
        <span v-if="isAdvancedUpload()"
          ><br />
          or drag it here</span
        >.
      </label>
    </div>
    <div>
      <span class="label">Title : </span>
      <input class="input_title" type="text" v-model="title" required />
    </div>
    <button @click="submitFile()">Submit</button>
  </form>
  <div v-else>
    <h5>{{ this.status }}</h5>
    <span>Uploading and transforming a video can be a lengthy process</span>
    <br />
    <span v-on:click="retry()">Retry</span>
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
  methods: {
    submitFile: function () {
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
      this.title = this.file.name;
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
      this.title = this.file.name;
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
  display: flex;
  flex-direction: column;
  row-gap: 10px;
  align-items: center;
}

#video_box {
  display: flex;
  float: right;
  height: 200px;
  width: 200px;
  border: 5px rgb(175, 175, 255) dashed;
  align-items: center;
  justify-content: center;
  &.is_dragged_over {
    border-style: solid;
    background-color: rgb(218, 218, 255);
  }
}

#drag_box {
  height: inherit;
  width: inherit;
  position: absolute;
  &.is_dragged_over {
    z-index: 2;
  }
}

.input_title {
  width: 200px;
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
</style>
