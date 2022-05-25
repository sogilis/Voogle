import { createRouter, createWebHashHistory } from "vue-router";
import HomePage from "../views/HomePage.vue";

const routes = [
  {
    path: "/",
    name: "HomePage",
    component: HomePage,
  },
  {
    path: "/watch/:id",
    name: "VideoPlayerPage",
    component: function () {
      return import("../views/VideoPlayerPage.vue");
    },
  },
  {
    path: "/upload",
    name: "UploadPage",
    component: function () {
      return import("../views/UploadPage.vue");
    },
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
