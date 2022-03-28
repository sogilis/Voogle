<template>
  <div class="login-container">
    <div v-if="cookies == undefined" class="login">
      <form @submit.prevent="login()">
        <input
          v-model="username"
          placeholder="Username"
          name="username"
          type="text"
        />
        <input
          placeholder="Password"
          v-model="password"
          name="password"
          type="password"
        />
        <button type="submit">Login</button>
      </form>
    </div>
    <div v-else class="logout">
      <button v-on:click="logout">Logout</button>
    </div>
  </div>
</template>

<script>
import cookies from "js-cookie";

export default {
  name: "Session",
  data: function () {
    return {
      username: null,
      password: null,
      cookies: undefined,
    };
  },
  mounted() {
    this.cookies = this.getCookies();
  },
  methods: {
    login: function () {
      cookies.set(
        "Authorization",
        "Basic " + btoa(this.username + ":" + this.password)
      );
      this.cookies = this.getCookies();
    },
    logout: function () {
      cookies.remove("Authorization");
      this.cookies = this.getCookies();
      this.$router.push("/");
    },
    getCookies: function () {
      this.$store.commit(
        "setLogState",
        cookies.get("Authorization") != undefined
      );
      return cookies.get("Authorization");
    },
  },
};
</script>
