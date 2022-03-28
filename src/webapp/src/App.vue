<template>
  <nav class="topnav">
    <router-link to="/">Home</router-link>
    <router-link to="/upload" v-show="isLogged">Upload</router-link>
    <Session></Session>
  </nav>
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

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: Arial, Helvetica, sans-serif;
}

h1 {
  font-size: 2em;
  margin: 30px;
  font-weight: bolder;
}

h2 {
  font-size: 1.5em;
  margin: 30px;
  font-weight: bolder;
}

nav {
  overflow: hidden;
  background-color: #e9e9e9;
  a {
    font-weight: bold;
    color: #2c3e50;
    float: left;
    display: block;
    text-align: center;
    padding: 14px 16px;
    text-decoration: none;
    font-size: 17px;
    transition: all 400ms;
    &:hover {
      background-color: #ddd;
      transform: scale(1.03);
    }

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}

.flex_center {
  display: flex;
  align-items: center;
  justify-content: center;
}

@media screen and (max-width: 600px) {
  nav a {
    float: none;
    display: block;
    text-align: left;
    width: 100%;
    margin: 0;
    padding: 14px;
  }
}
</style>
