<template>
  <div class="session">
    <div v-if="cookies == undefined">
      <form @submit.prevent="login()">
        <input
          v-model="username"
          placeholder="Username"
          name="username"
          type="text"
          class="session__input"
        />
        <input
          placeholder="Password"
          v-model="password"
          name="password"
          type="password"
          class="session__input"
        />
        <button type="submit" class="session__button">Login</button>
      </form>
    </div>
    <div v-else>
      <button v-on:click="logout" class="session__button">Logout</button>
    </div>
  </div>
</template>

<script>
import cookies from "js-cookie";

export default {
  name: "SessionForm",
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
        { sameSite: "lax" }
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

<style lang="scss">
.session {
  float: right;

  &__input {
    padding: 6px;
    margin-top: 8px;
    margin-left: 6px;
    font-size: 17px;
    border: none;
    width: 120px;
  }

  &__button {
    float: right;
    padding: 6px 10px;
    margin-top: 8px;
    margin-right: 16px;
    margin-left: 6px;
    background-color: #555;
    color: white;
    font-size: 17px;
    border: none;
    cursor: pointer;
    transition: background-color 400ms;

    &:hover {
      background-color: green;
    }
  }
}
@media screen and (max-width: 600px) {
  .session {
    float: none;
    &__input,
    &__button {
      float: none;
      display: block;
      text-align: left;
      width: 100%;
      margin: 0;
      padding: 14px;
    }
    &__input {
      border: 1px solid #ccc;
    }
  }
}
</style>
