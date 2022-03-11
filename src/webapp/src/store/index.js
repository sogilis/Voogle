import { createStore } from "vuex";

export default createStore({
  state: {
    isLogged: false,
  },
  getters: {},
  mutations: {
    setLogState(isLogged, newStatus) {
      this.state.isLogged = newStatus;
    },
  },
  actions: {},
  modules: {},
});
