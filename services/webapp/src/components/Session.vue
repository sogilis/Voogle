<template>
  <div>
    <div v-if="cookies == undefined" class="login">
      <h2>Login</h2>
      <form @submit.prevent="login">
        <p>
          <label for="username">Username</label>
          <input id="username" v-model="username" name="username" /><br />
          <label for="password">Password</label>
          <input id="password" v-model="password" name="password" />
        </p>
        <input type="submit" value="login" v-on:click="login" />
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
        "Basic " + btoa(this.username + ":" + this.password),
        { expires: 30 }
      );
      this.cookies = this.getCookies();
    },
    logout: function () {
      cookies.remove("Authorization");
      this.cookies = this.getCookies();
    },
    getCookies: function () {
      return cookies.get("Authorization");
    },
  },
};
</script>
