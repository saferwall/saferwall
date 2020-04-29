<template>
  <div class="container">
    <tree-view :data="data" :options="options" />
  </div>
</template>

<script>
export default {
  props: ["data", "sectionName"],
  data() {
    return {
      options: {
        maxDepth: 1,
        rootObjectKey: this.sectionName,
        modifiable: false,
        link: false,
        limitRenderDepth: false,
      },
    }
  },
  watch: {
    sectionName: function(val) {
      this.options.rootObjectKey = val
    },
  },
}
</script>

<style lang="scss">
.tree-view-item-value {
  color: #00d1be;
}
</style>
