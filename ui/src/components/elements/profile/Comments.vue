<template>
  <div v-if="active">
    <p id="no" v-if="!comments">No Comments</p>
    <CommentCard
      :comment="comment"
      v-for="(comment, index) in comments"
      :key="index"
    />
  </div>
</template>

<script>
import CommentCard from "./CommentCard"

export default {
  props: ["active", "comments"],
  components: {
    CommentCard,
  },
}
</script>

<style scoped>
#no {
  font-size: 25px;
  font-weight: 200;
  padding: 0.5em;
}
</style>
