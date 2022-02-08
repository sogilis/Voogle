import { expect } from "chai";
import { shallowMount } from "@vue/test-utils";
import Session from "@/components/Session.vue";

describe("Session.vue", () => {
  let component;

  beforeEach(() => {
    component = shallowMount(Session);
  });

  it("Renders input and button", () => {
    expect(component.find("button[type='submit']").exists()).to.be.true;
    expect(component.find("input[name='username']").exists()).to.be.true;
    expect(component.find("input[name='password']").exists()).to.be.true;
  });
});
