<template>
  <div>
    <div class="tabs" :class="type">
      <ul>
        <li
          v-for="tab in tabs"
          :key="tab.name"
          :class="{ 'is-active': tab.isActive }"
        >
          <a @click="selectTab(tab)" class="tab-title">
            <span class="icon is-small"><i :class="tab.$attrs.icon"></i></span>
            <span>{{ tab.name }}</span>
          </a>
        </li>
      </ul>
    </div>
    <div class="tab-details">
      <slot></slot>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      tabs: [],
    }
  },
  props: ["type", "url"],
  created() {
    this.tabs = this.$children
  },
  methods: {
    selectTab(selectedTab) {
      this.tabs.forEach((tab) => {
        tab.isActive =
          selectedTab.name === tab.name &&
          this.$emit("tabChanged", selectedTab.name)
      })
    },
  },
}
</script>
<style lang="scss" scoped>
@import "../../assets/scss/variables";

li {
  &.is-active {
    .tab-title {
      color: $primary-color;
    }
  }
}
</style>
