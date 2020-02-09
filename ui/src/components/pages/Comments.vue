<template>
  <div>
    <p id="noComments" v-if="!comments || comments.legth === 0">No Comments Available</p>
    <CommentCard
      v-for="comment in comments"
      :key="comment.id"
      :data="comment"
    />
    <TextEditor />
  </div>
</template>

<script>
import TextEditor from "../elements/comments/TextEditor"
import CommentCard from "../elements/comments/CommentCard"

import { mapGetters } from "vuex"

export default {
  components: {
    TextEditor,
    CommentCard,
  },
  computed: {
    ...mapGetters({
      comments: "getComments",
    }),
  },
}
</script>

<style>
#noComments {
  margin-bottom: 2em;
  font-size: 1.3em;
}
</style>
