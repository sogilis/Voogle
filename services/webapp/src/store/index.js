import { createStore } from 'vuex'

export default createStore({
  state: {
    clickCounter: 20
  },
  getters: {
  },
  mutations: {
    setClickCounter (state, value) {
      if (!isNaN(value)) {
        state.clickCounter = value
      }
    }
  },
  actions: {
  },
  modules: {
  }
})
