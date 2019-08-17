<template>
  <div :class="['notification', type]">
    <button class="delete" @click="closeNotif()"></button>
    <slot></slot>
  </div>
</template>
<script>
export default {
  props: ["type", "active"],
  methods: {
    closeNotif() {
      this.$emit("closeNotif")
    },
  },
}
</script>
<style lang="scss" scoped></style>
