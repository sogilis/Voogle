<template>
  <div>
    <h1>{{video.title}}</h1>
    <h2>{{video.description}}</h2>
  </div>
</template>

<script>
export default {
  props: {
    video: Object
  }
}
</script>
