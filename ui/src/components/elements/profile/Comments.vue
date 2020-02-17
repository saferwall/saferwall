<template>
  <div v-if="active">
    <CommentCard :comment="comment" v-for="(comment, index) in comments" :key="index"/>
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

<style></style>
