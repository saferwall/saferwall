<template>
  <div>
    <div class="tabs is-boxed">
      <ul>
        <li :class="{ 'is-active': activeTab === 0 }">
          <a @click="() => (activeTab = 0)">Editor</a>
        </li>
        <li :class="{ 'is-active': activeTab === 1 }">
          <a @click="() => (activeTab = 1)">Preview</a>
        </li>
      </ul>
    </div>
    <div class="editor" v-if="activeTab === 0">
      <quill-editor
        v-model="content"
        ref="myQuillEditor"
        :options="editorOption"
      >
      </quill-editor>
    </div>
    <div class="preview" v-if="activeTab === 1" v-html="content?content:'Nothing to Preview'"></div>
    <button class="button comment-btn" @click="sendComment">Comment</button>
  </div>
</template>

<script>
import "quill/dist/quill.core.css"
import "quill/dist/quill.snow.css"
import "quill/dist/quill.bubble.css"

import { quillEditor } from "vue-quill-editor"

const toolbarOptions = [
  ["bold", "italic", "underline", "strike"], // toggled buttons
  ["blockquote", "code-block"],
  [{ header: 1 }, { header: 2 }], // custom button values
  [{ list: "ordered" }, { list: "bullet" }],
  [{ script: "sub" }, { script: "super" }], // superscript/subscript
]

export default {
  components: {
    quillEditor,
  },
  data() {
    return {
      activeTab: 0,
      content: "",
      editorOption: {
        modules: {
          toolbar: toolbarOptions,
        },
      },
    }
  },
  methods: {
    sendComment: function() {
      if (this.content !== "") {
        this.$http.post(this.$api_endpoints.FILES+this.$store.getters.getHashContext+"/comments/", {
          body:this.content
        })
        .then(()=>{
          this.$store.dispatch('updateComments')
          this.content = ""
        })
        .catch(()=>{
          this.$awn.alert('An Error Occured, Try Again')
        })
      }
    },
  },
}
</script>

<style lang="scss">
.preview {
  border: 1px solid #ccc;
  padding: 15px;

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
.comment-btn {
  margin-top: 1rem;
}
</style>
