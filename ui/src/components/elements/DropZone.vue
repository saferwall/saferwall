<template>
  <div
    class="tile is-child box dropzone has-padding"
    @dragover.prevent="onDragOver"
    @dragleave.prevent="onDragLeave"
    @drop.prevent="onDrop"
    :class="{ highlight: highlight }"
  >
    <input
      ref="fileInput"
      @change.prevent="onFileChange"
      type="file"
      class="file-input"
    />
    <div class="icon"><i class="ion-ios-cloud-upload" /></div>
    <h2 class="title is-4">Drag &amp; Drop files here</h2>

    <h2 class="title is-4 ">Or</h2>
    <button class="btn" @click="openFileDialog" type="submit">
      Browse files
    </button>
  </div>
</template>

<script>
export default {
  data() {
    return { highlight: false };
  },
  methods: {
    onDragOver() {
      this.highlight = true;
    },
    onDragLeave() {
      this.highlight = false;
    },
    onDrop(e) {
      const file = e.dataTransfer.files.item(0);
      this.$emit("fileAdded", file);
      this.highlight = false;
    },
    openFileDialog() {
      this.$refs.fileInput.click();
    },
    onFileChange(e) {
      const files = e.target.files;
      const fileIsSelected = files.length > 0;

      if (fileIsSelected) {
        const file = e.target.files.item(0);
        this.$emit("fileAdded", file);
        this.highlight = false;
      } else {
        alert("Please select a file to upload");
      }
    }
  }
};
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
