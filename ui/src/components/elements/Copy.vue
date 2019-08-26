<template>
  <div class="copy" :class="{ copied: copied }" @click="copy($event)">
    Copy
  </div>
</template>
<script>
export default {
  props: ["content"],
  data() {
    return {
      copied: false,
    }
  },
  methods: {
    copy(e) {
      this.copied = true
      this.$clipboard.copy(this.content)
      setTimeout(() => {
        this.copied = false
      }, 500)
    },
  },
}
</script>
<style lang="scss" scoped>
@import "../../assets/scss/variables";
@keyframes copying {
  from {
    transform: translate(50%, -50%) scale(1);
    opacity: 1;
  }
  to {
    transform: translate(50%, -50%) scale(1.2);
    opacity: 0;
  }
}
.copy {
  position: absolute;
  top: 50%;
  right: 50%;
  transform: translate(50%, -50%);
  display: inline-block;
  color: #fff;
  font-weight: 500;
  font-size: 12px;
  font-weight: 500;
  padding: 0 3px;
  border-radius: 3px;
  cursor: pointer;
  background-color: $primary-color;
  z-index: 999;

  &.copied {
    animation: copying 0.5s;
  }
}
</style>
