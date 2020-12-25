<template>
  <div class="foot">
    <div class="icons">
      <i
        @click="openSocial('github')"
        class="icon mdi mdi-24px mdi-github-face"
        style="color:black;"
      ></i>
      <i
        @click="openSocial('twitter')"
        class="icon mdi mdi-24px mdi-twitter"
        style="color:#1DA1F2;"
      ></i>
    </div>
    <p>
      Made with <span class="heart">&#10084;</span> by <a href="/">SaferWall</a>
    </p>
  </div>
</template>

<script>
export default {
  methods: {
    openSocial(social) {
      switch (social) {
        case "github":
          window.open("https://github.com/saferwall/saferwall", "_bank")
          break
        case "twitter":
          window.open("https://twitter.com/saferwall", "_bank")
          break
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.foot {
  text-align: center;
  margin-top: 2rem;
}
.icons {
  & .icon:first-child {
    margin-right: 1em;
  }
  .icon {
    cursor: pointer;
  }
}
.heart {
  color: red;
  font-size: 1.3em;
  text-align: center;
  vertical-align: sub;
}
</style>
