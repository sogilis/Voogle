<template>
  <div id="nav" class="topnav">
    <router-link to="/">Home</router-link>
    <router-link to="/upload" v-show="isLogged">Upload</router-link>
    <Session></Session>
  </div>
  <router-view />
</template>

<script>
import Session from "@/components/Session";
export default {
  components: { Session },
  computed: {
    isLogged: function () {
      return this.$store.state.isLogged;
    },
  },
};
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: Arial, Helvetica, sans-serif;
}

.topnav {
  overflow: hidden;
  background-color: #e9e9e9;
}

.topnav a {
  float: left;
  display: block;
  color: black;
  text-align: center;
  padding: 14px 16px;
  text-decoration: none;
  font-size: 17px;
}

.topnav a:hover {
  background-color: #ddd;
  color: black;
}

.topnav a.active {
  background-color: #2196f3;
  color: white;
}

.topnav .login-container {
  float: right;
}

.topnav input[type="text"],
.topnav input[type="password"] {
  padding: 6px;
  margin-top: 8px;
  margin-left: 6px;
  font-size: 17px;
  border: none;
  width: 120px;
}

.topnav .login-container button {
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
}

.topnav .login-container button:hover {
  background-color: green;
}

@media screen and (max-width: 600px) {
  .topnav .login-container {
    float: none;
  }
  .topnav a,
  .topnav input[type="text"],
  .topnav input[type="password"],
  .topnav .login-container button {
    float: none;
    display: block;
    text-align: left;
    width: 100%;
    margin: 0;
    padding: 14px;
  }
  .topnav input[type="text"],
  .topnav input[type="password"] {
    border: 1px solid #ccc;
  }
}
</style>
