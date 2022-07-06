<template>
  <div v-if="this.withSort">
    <label for="attribute">Sort by :</label><br />
    <select
      name="attribute"
      @change="selectChange($event.target.value, this.ascending, this.status)"
      :value="this.attribute"
    >
      <option value="title">Title</option>
      <option value="upload_date">Upload Date</option>
    </select>
    <select
      name="ascending"
      @change="selectChange(this.attribute, $event.target.value, this.status)"
      :value="this.ascending"
    >
      <option value="true">Ascending</option>
      <option value="false">Descending</option>
    </select>
    <br />
    <label for="status">Show :</label><br />
    <select
      name="status"
      @change="
        selectChange(this.attribute, this.ascending, $event.target.value)
      "
      :value="this.status"
    >
      <option value="Complete">Uploaded</option>
      <option value="Archive">Archived</option>
    </select>
  </div>
  <div class="PageNavigation">
    <button
      class="button PageNavigation__button"
      @click="pageChange('first')"
      :disabled="this.is_first"
    >
      <i class="fa-solid fa-backward-fast"></i>
    </button>
    <button
      class="button PageNavigation__button"
      @click="pageChange('previous')"
      :disabled="this.is_first"
    >
      <i class="fa-solid fa-caret-left"></i>
    </button>
    {{ this.page }}
    <button
      class="button PageNavigation__button"
      @click="pageChange('next')"
      :disabled="this.is_last"
    >
      <i class="fa-solid fa-caret-right"></i>
    </button>
    <button
      class="button PageNavigation__button"
      @click="pageChange('last')"
      :disabled="this.is_last"
    >
      <i class="fa-solid fa-forward-fast"></i>
    </button>
  </div>
</template>

<script>
export default {
  name: "PageNavigation",
  props: {
    page: Number,
    is_last: Boolean,
    is_first: Boolean,
    attribute: String,
    ascending: Boolean,
    status: String,
    withSort: Boolean,
  },
  emits: ["selectChange", "pageChange"],
  methods: {
    selectChange: function (attr, asc, stat) {
      asc = asc == "true";
      this.$emit("selectChange", {
        attribute: attr,
        ascending: asc,
        status: stat,
      });
    },
    pageChange: function (i) {
      this.$emit("pageChange", { page: i });
    },
  },
};
</script>

<style scoped lang="scss">
.PageNavigation {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  column-gap: 1em;
  margin: 1em;
  font-size: 1.5em;

  &__button {
    background-color: rgb(241, 241, 241);
  }
}
</style>
