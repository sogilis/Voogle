<template>
  <div>
    <button @click="pageChange('first')">First</button>
    <button @click="pageChange('previous')">Previous</button>
    {{ this.page }}
    <button @click="pageChange('next')" :disabled="this.is_last">Next</button>
    <button @click="pageChange('last')" :disabled="this.is_last">Last</button>
  </div>
</template>

<script>
export default {
  name: "Navigation",
  props: {
    page: Number,
    is_last: Boolean,
  },
  methods: {
    pageChange: function (i) {
      this.$emit("pageChange", { page: i });
    },
  },
};
</script>
