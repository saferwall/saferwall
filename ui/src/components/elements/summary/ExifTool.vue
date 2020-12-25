<template>
  <div class="box">
    <h4 class="title">ExifTool File Metadata</h4>
    <div v-for="(i, index) in summaryData.exif" :key="index" class="data-data">
      <strong class="data-label exif">
        {{ index.replace(/[A-Z]/g, (match) => ` ${match}`) }}
      </strong>
      <span class="data-value">
        <span class="value-text">{{ i }}</span>

        <copy :content="i"></copy>
      </span>
    </div>
  </div>
</template>

<script>
import Copy from "@/components/elements/Copy"

export default {
  props: ["summaryData"],
  components: {
    copy: Copy,
  },
}
</script>

<style></style>
