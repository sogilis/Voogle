import { expect } from "chai";
import { mount } from "@vue/test-utils";
import SessionForm from "@/components/SessionForm.vue";
import { createStore } from "vuex";
import { createRouter, createWebHashHistory } from "vue-router";

const store = createStore({
  mutations: {
    setLogState(state, newStatus) {
      state.isLogged = newStatus;
    },
  },
});

const router = createRouter({
  history: createWebHashHistory(),
  routes: [],
});

describe("SessionForm.vue", () => {
  const component = mount(SessionForm, {
    global: {
      plugins: [[store], [router]],
    },
  });

  it("Renders input and button", () => {
    expect(component.find("button[type='submit']").exists()).to.be.true;
    expect(component.find("input[name='username']").exists()).to.be.true;
    expect(component.find("input[name='password']").exists()).to.be.true;
  });
});
