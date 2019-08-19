<template>
  <div>
    <header></header>
    <main>
      <router-view></router-view>
    </main>
  </div>
</template>
<script>
export default {};
</script>
<style lang="scss" scoped>
header {
  --center: 22vh;
  --ascend: 6.613vh;

  height: calc(var(--center) + var(--ascend));
  background: #e0e0e0;

  clip-path: polygon(
    0 0,
    100% 0,
    100% calc(var(--center) - var(--ascend)),
    0 calc(var(--center) + var(--ascend))
  );
}
.container {
  transition: 0.8s cubic-bezier(0.2, 0.8, 0.2, 1);
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: 1fr 1fr;
  grid-gap: 0.5em;
  align-items: center;
  padding: 2em;
  font-size: 0.7rem;
}
@media (min-width: 640px) {
  .container {
    font-size: 1rem;
    max-width: 100%;
    justify-items: center;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr;
  }
}
</style>
