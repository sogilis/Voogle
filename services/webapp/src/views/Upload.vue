<template>
  <h2>UPLOAD</h2>
  <div v-if="this.status === 'None'">
    <label>
      Title
      <input type="text" v-model="title" /> </label
    ><br />
    <label
      >File
      <input type="file" ref="file" v-on:change="handleFileUpload()" />
    </label>
    <br />
    <button v-on:click="submitFile()">Submit</button>
  </div>
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
        .post(
          process.env.VUE_APP_API_ADDR + "api/v1/videos/upload",
          formData,
          {
            headers: {
              "Content-Type": "multipart/form-data",
              Authorization: cookies.get("Authorization"),
            },
          }
        )
        .then((res) => {
          if (res.status !== 200) {
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
  },
};
</script>

<style scoped></style>
