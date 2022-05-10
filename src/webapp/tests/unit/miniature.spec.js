import { expect } from "chai";
import { shallowMount } from "@vue/test-utils";
import Miniature from "@/components/Miniature.vue";

describe("Miniature.vue", async () => {
  const wrapper = shallowMount(Miniature);
  const testtitle = "title";

  await wrapper.setProps({ title: testtitle });
  await wrapper.setProps({ enable_deletion: false });

  it("Renders title", () => {
    console.log(wrapper.text());
    expect(wrapper.text()).to.include(testtitle);
  });

  it("Has no delete button", () => {
    expect(wrapper.find(`button.miniature__delete-button`).exists()).to.be
      .false;
  });

  describe("Enable Deletion", async () => {
    before(async () => {
      await wrapper.setProps({ enable_deletion: true });
    });
    it("Has a delete button", () => {
      expect(wrapper.find(`button.miniature__delete-button`).exists()).to.be
        .true;
    });
  });
});
