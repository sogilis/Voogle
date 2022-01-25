import {createRouter, createWebHashHistory} from "vue-router";
import Home from "../views/Home.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/watch/:id",
    name: "Watch",
    component: function () {
      return import("../views/Watch.vue");
    },
  },
  {
    path: "/upload",
    name: "Upload",
    component: function () {
      return import("../views/Upload.vue");
    },
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
