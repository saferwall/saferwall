<template>
  <div class="columns" style="margin-top:100px">
    <div class="column is-8 is-offset-2">
      <tabs type="is-boxed">
        <tab name="File" icon="ion-android-folder-open" :selected="true">
          <DropZone @fileAdded="onFileAdded" /><transition
            name="slide"
            mode="out-in"
          >
            <notification
              type="is-danger"
              @closeNotif="close()"
              v-if="notifActive"
            >
              {{ notificationError }}
            </notification>
          </transition>
        </tab>
        <tab name="Url" icon="ion-ios-world-outline">
          <form class="tile is-child box">
            <h1>This part is about our culture</h1>

            <p class="is-centered" style="margin-top:10px;">
              <small>
                By using Saferwall you consent to our
                <router-link to="/">Terms of Service</router-link> and
                <router-link to="">Privacy Policy</router-link> and allow us to
                share your submission with the security community.
                <router-link to="">Learn more.</router-link>
              </small>
            </p>
          </form>
        </tab>
      </tabs>
    </div>
  </div>
</template>
<script>
import Tabs from "@/components/elements/Tabs";
import Tab from "@/components/elements/Tab";
import Notification from "@/components/elements/Notification";
import { sha256 } from "js-sha256";
import axios from "axios";
import Scanning from "@/components/pages/Scanning";
import DropZone from "@/components/elements/DropZone";

export default {
  data() {
    return {
      notificationError: "",
      notifActive: false,
      filename: ""
    };
  },
  components: {
    tabs: Tabs,
    tab: Tab,
    notification: Notification,
    DropZone
  },
  methods: {
    close() {
      this.notifActive = false;
    },
    onFileAdded(file) {
      if (!file) {
        return;
      }
      // check if size exceeds 64mb
      if (file.size > 64000000) {
        this.notifActive = true;
        this.notificationError = "File size exceeds 64MB !";
        return;
      }
      const reader = new FileReader();
      reader.onload = loadEvent => {
        // file has been read successfully
        const fileBuffer = loadEvent.target.result;
        crypto.subtle
          .digest("SHA-256", fileBuffer)
          .then(hashBuffer => {
            const hashArray = new Uint8Array(hashBuffer);
            let hashHex = "";
            for (let i = 0; i < hashArray.byteLength; i++) {
              let hex = new Number(hashArray[i]).toString("16");
              if (hex.length == 1) {
                hex = "0" + hex;
              }
              hashHex += hex;
            }
            // hash hexadecimal has been calculated successfully
            axios
              .get(`/api/v1/files/${hashHex}`)
              .then(response => {
                this.$router.push(`summary/${hashHex}`);
              })
              .catch(
                // upload the file to the db
                // perform scanning 
                // show scanning progress on /scanning
              );
          })
          .catch(error => {
            this.notifActive = true;
            this.notificationError =
              "Sorry, we couldn't upload the file. Please, try again!";
          });
      };
      reader.readAsArrayBuffer(file);
    }
  }
};
</script>
<style lang="scss" scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}
.slide-enter,
.slide-leave-to {
  height: 0;
  padding: 0;
  overflow: hidden;
  opacity: 0;
}

.tile {
  margin-top: -24px !important;
  border-radius: 0 0 4px 4px !important;
  margin-left: 1px !important;
}
.file-container {
  display: inline-block;
  margin: auto;
  position: relative;
  text-align: center;
  cursor: pointer;
  padding: 50px 0;

  input[type="file"] {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    opacity: 0;
    cursor: pointer;
    z-index: 9;
  }
}
a.button {
  display: inline-block;
  margin-top: 10px;
}
</style>
