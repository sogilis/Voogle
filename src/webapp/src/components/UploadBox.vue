<template>
  <div id="upload_box_id">
    <div id="back" class="uploadbox" :class="{ set: videoIsSet }">
      <input
        v-bind:id="refto"
        class="uploadbox__input"
        type="file"
        v-bind:ref="refto"
        @change="handleFileSelect()"
        v-bind:accept="accepting"
      />
      <label class="uploadbox__text" v-bind:for="refto" v-if="!videoIsSet">
        <strong>Choose a file</strong>
        <span v-if="dragEnabled"
          ><br />
          or drag it here</span
        >.
      </label>
      <!-- Div handling drag-events when supported -->
      <div
        id="front"
        class="uploadbox__dragbox"
        v-if="dragEnabled"
        @drop.prevent.stop="handleDrop"
        @dragleave.prevent.stop="dragOnBox(false)"
        @dragend.prevent.stop="dragOnBox(false)"
        @dragenter.prevent.stop="dragOnBox(true)"
        @dragover.prevent.stop="dragOnBox(true)"
        @drag.prevent.stop=""
        @dragstart.prevent.stop=""
      ></div>
      <div class="uploadbox__preview" v-if="videoIsSet">
        <img
          class="uploadbox__preview-img"
          src="../assets/uxwing-camera.png"
          alt="Video Set"
        />
        <div class="uploadbox__preview-title">
          <p>{{ title }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "UploadBox",
  el: "#upload_box_id",
  data: function () {
    return {};
  },
  props: {
    title: String,
    accepting: String,
    refto: String,
  },
  computed: {
    videoIsSet: function () {
      return this.title;
    },
    dragEnabled: function () {
      var div = document.createElement("div");
      return (
        ("draggable" in div || ("ondragstart" in div && "ondrop" in div)) &&
        "FormData" in window &&
        "FileReader" in window
      );
    },
  },
  methods: {
    handleDrop: function (e) {
      this.$emit("sendFile", { file: e.dataTransfer.files[0] });
    },
    handleFileSelect() {
      this.$emit("sendFile", { file: this.$refs[this.refto].files[0] });
    },
    dragOnBox: function (bool) {
      var back = this.$el.querySelector("#back");
      var front = this.$el.querySelector("#front");

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
.uploadbox {
  display: flex;
  height: 200px;
  width: 200px;
  border: 5px rgb(175, 175, 255) dashed;
  border-radius: 10px;
  overflow: hidden;
  align-items: center;
  justify-content: center;

  &__input {
    display: none;
  }

  &.set {
    border: 2px rgb(175, 175, 255) solid;
  }

  &__text {
    z-index: 1;
    &:hover strong {
      text-decoration: underline;
    }
  }

  &__dragbox {
    position: absolute;
    height: inherit;
    width: inherit;
    &.is_dragged_over {
      z-index: 2;
    }
  }

  &.is_dragged_over {
    border-style: solid;
    background-color: rgb(218, 218, 255);
  }

  &__preview {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    width: 100%;

    &-title {
      position: absolute;
      bottom: 0px;
      padding: 5px;
      width: 100%;
      max-height: 2rem;
      text-align: center;
      background-color: black;
      color: white;
    }

    &-img {
      width: 70%;
    }
  }
}
</style>
