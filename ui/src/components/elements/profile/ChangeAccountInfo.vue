<template>
  <div class="columns page">
    <div class="column" style="flex-grow:0.6;">
      <div class="form-group">
        <label class="form__label">Name</label>
        <input
          class="input"
          v-model.trim="$v.name.$model"
          type="text"
          placeholder="Real name"
        />
      </div>
      <div class="form-group" :class="{ 'form-group--error': $v.bio.$error }">
        <label class="form__label">Bio</label>
        <textarea
          class="textarea has-fixed-size"
          v-model.trim="$v.bio.$model"
          placeholder="Small Bio"
          rows="5"
        />
        <div class="error" v-if="!$v.bio.max && $v.bio.$dirty">
          Bio should contain at most 200 characters
        </div>
      </div>
      <div class="form-group">
        <label class="form__label">Member Since</label>
        <input class="input" :value="member_since" type="text" readonly />
      </div>
      <div class="form-group">
        <label class="form__label">Location</label>
        <input
          class="input"
          v-model.trim="$v.location.$model"
          type="text"
          placeholder="location"
        />
      </div>
      <div class="form-group">
        <label class="form__label">URL</label>
        <input
          class="input"
          v-model.trim="$v.URL.$model"
          type="text"
          placeholder="url"
        />
      </div>
      <button
        class="button is-primary"
        :disabled="!tmp_avatar && !$v.$anyDirty"
        @click="submit"
      >
        {{ this.submitLabel }}
      </button>
    </div>
    <div class="column">
      <figure class="image">
        <img
          v-if="!tmp_avatar"
          :src="
            userData.avatarBase64
              ? 'data:image/png;base64,' + userData.avatarBase64
              : ''
          "
        />
        <img v-if="tmp_avatar" :src="tmp_avatar" />

        <button class="button is-small image-overlay" @click="selectImage">
          <i class="icon mdi mdi-pencil"></i>
          &nbsp;Modify
        </button>
      </figure>
      <input
        type="file"
        style="display:none;"
        ref="imageInput"
        accept="image/png"
        @change="imageAdded"
      />
    </div>
  </div>
</template>

<script>
import { maxLength } from "vuelidate/lib/validators"
import { mapGetters } from "vuex"

export default {
  data() {
    return {
      name: "",
      bio: "",
      location: "",
      member_since: "",
      URL: "",
      tmp_avatar: null,
      file: null,
    }
  },
  computed: {
    ...mapGetters({
      userData: "getUserData",
    }),
    submitLabel: function() {
      if (this.$v.$anyDirty) return "Submit"
      else if (this.tmp_avatar) return "Change Avatar"
      else return "Submit"
    },
  },
  methods: {
    loadData() {
      if (
        !this.userData ||
        (Object.entries(this.userData).length === 0 &&
          this.userData.constructor === Object)
      )
        return

      this.name = this.userData.name?this.userData.name:this.userData.username
      this.bio = this.userData.bio
      this.location = this.userData.location
      this.URL = this.userData.url
      this.member_since = this.userData.member_since.substring(0, 10)
    },
    submit: function() {
      var updated = false
      if (!this.$v.$invalid && this.$v.$anyDirty) {
        this.submitInfo()
        updated = true
      }
      if (this.tmp_avatar !== null) {
        this.submitImage()
        updated = true
      }
      if (updated)
        this.$store.dispatch("updateUserData", this.userData.username)
    },
    submitInfo: function() {
      var data = {}
      if (this.name) data.name = this.name
      if (this.bio) data.bio = this.bio
      if (this.location) data.location = this.location
      if (this.URL) data.url = this.URL

      this.$http
        .put(this.$api_endpoints.USERS + this.$store.getters.getUsername, data)
        .then(() => {
          this.$awn.success("Information Changed Successfully")
          this.trackSuccess()
          this.clear()
        })
        .catch(() => {
          this.$awn.alert("A problem occured, try again")
        })
    },
    submitImage: function() {
      const reader = new FileReader()
      reader.onload = (loadEvent) => {
        // file has been read successfully
        const formData = new FormData()
        formData.append("file", this.file)
        this.$http
          .put(
            this.$api_endpoints.USERS + this.userData.username + "/avatar",
            formData,
            {
              headers: {
                "Content-Type": "multipart/form-data",
              },
            },
          )
          .then(() => {
            this.$awn.success("Avatar Changed Successfully")
          })
          .catch(() => {
            this.$awn.alert("An Error Occured, Try Again")
          })
      }
      reader.readAsArrayBuffer(this.file)
    },
    trackSuccess() {
      this.$gtag.event("Account_information_change", {
        event_category: "Information_Change",
        event_label: "Info Changed",
        value: 1,
      })
    },
    clear() {
      this.$v.$reset()
    },
    selectImage() {
      this.$refs.imageInput.click()
    },
    imageAdded(e) {
      this.file = e.srcElement.files[0]

      if (this.file.size > 400000) {
        this.$awn.alert("Image too big!")
        return
      }

      if (!this.file) {
        this.$awn.alert("File cannot be read!")
        return
      }
      this.tmp_avatar = URL.createObjectURL(this.file)
    },
  },
  mounted() {
    this.loadData()
  },
  validations: {
    bio: {
      max: maxLength(255),
    },
    name: {},
    URL: {},
    location: {},
  },
}
</script>

<style lang="scss" scoped>
.page {
  margin-left: 5em !important;
  .form-group {
    display: grid;
    margin-bottom: 1em;

    .form__label {
      font-weight: 600;
    }

    .input {
      width: 25em;
    }

    .textarea {
      width: 25em;
      min-width: 25em;
    }
  }

  .error {
    color: rgb(243, 54, 21);
  }

  .button {
    margin-top: 0.5em;
    width: 7em;
  }

  .form-group--error {
    .input {
      border-color: rgb(243, 54, 21);
    }
  }

  .image {
    width: 200px;
    &:hover {
      .image-overlay {
        display: inline;
      }
    }
  }

  .image-overlay {
    position: absolute;
    bottom: 0;
    display: none;
  }
}
</style>
