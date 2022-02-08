import { expect } from "chai";
import { shallowMount } from "@vue/test-utils";
import Miniature from "@/components/Miniature.vue";

describe("Miniature.vue", () => {
  it("Renders props.title when passed", () => {
    const title = "title test";
    const wrapper = shallowMount(Miniature, {
      props: { title },
    });
    expect(wrapper.text()).to.include(title);
  });
});
