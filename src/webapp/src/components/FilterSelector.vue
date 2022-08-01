<template>
  <div class="filterSelector">
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
    {{ this.errmsg }}
  </div>
</template>

<script>
import axios from "axios";
import cookies from "js-cookie";
export default {
  name: "FilterSelector",
  data: function () {
    return {
      filterlist: [],
      errmsg: "",
    };
  },
  mounted() {
    axios
      .get(process.env.VUE_APP_API_ADDR + `api/v1/videos/transformer/list`, {
        headers: {
          Authorization: cookies.get("Authorization"),
        },
      })
      .then((response) => {
        var services = response.data["services"];
        services.forEach((service) => {
          this.filterlist.push({
            label: service["name"],
            name: service["name"],
            value: false,
          });
        });
      })
      .catch((error) => {
        this.errmsg = error;
      });
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

<style scoped lang="scss">
.filterSelector {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  column-gap: 4rem;
}
</style>
