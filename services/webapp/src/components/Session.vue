<template>
  <div class="login">
    <h2>Login</h2>
    <form @submit.prevent="login">
      <p>
        <label for="username">Username</label>
        <input id="username" v-model="username" name="username"><br>
        <label for="password">Password</label>
        <input id="password" v-model="password" name="password">
      </p>
      <input type="submit" value="login" v-on:click="login">
    </form>
  </div>
</template>

<script>
import cookies from 'js-cookie';

export default {
  name: "Session",
  data: function () {
    return {
      username: null,
      password: null,
    }
  },
  methods:
    {
      login: function () {
        cookies.set('Authorization', "Basic " + btoa(this.username + ":" + this.password), {expires: 30})
      }
    }
}
</script>

<style scoped>

</style>
