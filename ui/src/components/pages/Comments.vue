<template>
  <div>
    <p id="noComments" v-if="!comments || comments.legth === 0">
      No Comments Available
    </p>
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

<style lang="scss">
#noComments {
  margin-bottom: 2em;
  font-size: 1.3em;
}
.comment_body {
  h1 {
    font-size: 2rem;
  }
  h2 {
    font-size: 1.5rem;
  }
  blockquote {
    margin-left: 32px;
    border-left: 4px solid #ccc;
    padding-left: 8px;
  }
  .ql-syntax {
    background-color: #23241f;
    color: #f8f8f2;
    overflow: visible;
    white-space: pre-wrap;
    margin-bottom: 5px;
    margin-top: 5px;
    padding: 5px 10px;
  }
  ol {
    padding-left: 1.5em;
  }
  ul {
    padding-left: 1.5em;
  }
}
</style>
