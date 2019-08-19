<template>
  <div class="container">
    <img
      class="illustration"
      src="@/assets/imgs/forbidden.svg"
      alt="403 illustration"
    />
    <div class="content">
      <h1 class="heading">Error: 403 Forbidden</h1>
      <p class="error-description">
        It seems like the page you’re looking for was deleted, never existed, or
        was mysteriously lost somewhere on the Internet.
      </p>
      <router-link to="/" tag="button" class="head-home">Head back home</router-link>
    </div>
  </div>
</template>
<script>
export default {};
</script>
<style lang="scss" scoped>
.error-description {
  font-size: 1.6em;
}
.illustration {
  width: 600px;
  padding: 1em;
}
.content {
  height: 100%;
  padding: 1em;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
}
.heading {
  font-size: 2,4em;
  font-weight: 600;
}
.head-home {
  margin: 2em 0;
  width: max-content;
  cursor: pointer;
  border-radius: 0.25rem;
  font-size: inherit;
  padding: 0.7em;
  font-weight: 600;
  color: white;
  background-color: #e7501d;
  border: none;
}
</style>
