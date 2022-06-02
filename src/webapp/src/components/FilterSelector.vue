<template>
  <div>
    <div class="field" v-for="(filter, index) in filterlist" :key="index">
      <input
        :id="filter[`name`]"
        type="checkbox"
        class="switch"
        checked="checked"
        v-model="this.filterlist[index][`value`]"
        @change="updateFilterList"
      />
      <label :for="filter[`name`]">{{ filter["label"] }}</label>
    </div>
  </div>
</template>

<script>
export default {
  name: "FilterSelector",
  data: function () {
    return {
      filterlist: [
        { label: "Black&White", name: "gray", value: false },
        { label: "Vertical flip", name: "flip", value: false },
      ],
    };
  },
  methods: {
    updateFilterList: function () {
      var filters = [];
      for (var index in this.filterlist) {
        var filter = this.filterlist[index];
        if (filter[`value`]) {
          filters.push(filter[`name`]);
        }
      }
      this.$emit("filterListUpdate", { filterList: filters });
    },
  },
};
</script>
