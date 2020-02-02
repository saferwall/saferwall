<template>
  <div class="parentContainer">
    <router-link to="/">
      <img
        sizes="(max-width: 150px) 100vw, 150px"
        srcset="
          ../assets/imgs/saferwall_jkdgq6_c_scale_w_200.png 200w,
          ../assets/imgs/saferwall_jkdgq6_c_scale_w_397.png 397w,
          ../assets/imgs/saferwall_jkdgq6_c_scale_w_549.png 549w,
          ../assets/imgs/saferwall_jkdgq6_c_scale_w_800.png 800w
        "
        src="saferwall_jkdgq6_c_scale_w_800.png"
        alt="saferwall logo"
    /></router-link>
    <router-view></router-view>
  </div>
</template>
<script>
export default {}
</script>

<style scoped>
.parentContainer {
  display: grid;
  justify-content: center;
  justify-items: center;
}
</style>
<style lang="scss" scoped></style>
