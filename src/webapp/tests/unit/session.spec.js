import { expect } from "chai";
import { mount } from "@vue/test-utils";
import Session from "@/components/Session.vue";
import { createStore } from "vuex";

const store = createStore({
  mutations: {
    setLogState(state, newStatus) {
      state.isLogged = newStatus;
    },
  },
});

describe("Session.vue", () => {
  const component = mount(Session, {
    global: {
      plugins: [store],
    },
  });

  it("Renders input and button", () => {
    expect(component.find("button[type='submit']").exists()).to.be.true;
    expect(component.find("input[name='username']").exists()).to.be.true;
    expect(component.find("input[name='password']").exists()).to.be.true;
  });
});
