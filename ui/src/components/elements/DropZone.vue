<template>
  <form
    method="post"
    :action="this.$api_endpoints.FILES"
    enctype="multipart/form-data"
    class="tile is-child box dropzone has-padding"
    @dragover.prevent="onDragOver"
    @dragleave.prevent="onDragLeave"
    @drop.prevent="onDrop"
    @submit.prevent
    :class="{ highlight: highlight, disabled: !enabled }"
  >
    <input
      ref="fileInput"
      @change.prevent="onFileChange"
      type="file"
      name="file"
      class="file-input"
    />
    <div class="icon"><i class="ion-ios-cloud-upload" /></div>
    <h2 class="title is-4">Drag &amp; Drop files here</h2>

    <h2 class="title is-4 ">Or</h2>
    <button class="btn" @click="openFileDialog" type="submit">
      Browse files
    </button>
    <p class="is-centered" style="margin-top:1.5em;">
      <small>
        By using Saferwall you consent to our
        <a href="https://about.saferwall.com/tos" target="_blank">Terms of Service</a>
        and <router-link to="">Privacy Policy</router-link> and allow us to
        share your submission with the security community.
        <router-link to="">Learn more.</router-link>
      </small>
    </p>
  </form>
</template>

<script>
export default {
  props: {
    enabled: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return { highlight: false }
  },
  methods: {
    onDragOver() {
      if (!this.enabled) return
      this.highlight = true
    },
    onDragLeave() {
      this.highlight = false
    },
    onDrop(e) {
      if (!this.enabled) return
      const files = e.dataTransfer.files
      const fileIsSelected = files.length > 0

      if (fileIsSelected) {
        const file = e.dataTransfer.files.item(0)
        this.$emit("fileAdded", file)
        this.highlight = false
      } else {
        alert("Please select a file to upload")
      }
    },
    openFileDialog() {
      if (!this.enabled) return
      this.$refs.fileInput.click()
    },
    onFileChange(e) {
      const files = e.target.files
      const fileIsSelected = files.length > 0

      if (fileIsSelected) {
        const file = e.target.files.item(0)
        this.$emit("fileAdded", file)
        this.highlight = false
      } else {
        alert("Please select a file to upload")
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.dropzone {
  height: 300px;
  text-align: center;
  transition: all cubic-bezier(0.39, 0.58, 0.57, 1) 0.1s;
}
.icon {
  width: 100%;
  height: 7rem;
}
.icon > i {
  font-size: 7rem;
  color: #02a678;
}

.title {
  width: 100%;
  text-align: center;
}

.disabled {
  opacity: 0.2;
  .btn {
    cursor: no-drop;
  }
}

.btn {
  cursor: pointer;
  border-radius: 0.25rem;
  font-size: inherit;
  padding: 0.7em;
  font-weight: 600;
  color: white;
  background-color: #e7501d;
  border: none;
}
.has-padding {
  padding: 3rem;
}
.highlight {
  outline: 2px solid #02a678;
  box-shadow: 4px 4px 0 0 #02a678;
}
.file-input {
  display: none;
}
</style>
