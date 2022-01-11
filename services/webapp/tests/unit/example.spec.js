import { expect } from 'chai'
import { shallowMount } from '@vue/test-utils'
import Gallery from '@/components/Gallery.vue'

describe('Gallery.vue', () => {
  it('renders props.msg when passed', () => {
    const msg = 'new message'
    const wrapper = shallowMount(Gallery, {
      props: { msg }
    })
    expect(wrapper.text()).to.include(msg)
  })
})
