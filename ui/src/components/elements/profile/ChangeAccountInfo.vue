<template>
  <div class="page">
    <div class="form-group" :class="{ 'form-group--error': $v.name.$error }">
      <label class="form__label">Name</label>
      <input class="input" v-model.trim="$v.name.$model" type="text" placeholder="Real Name"/>
      <div class="error" v-if="!$v.name.required && $v.name.$dirty">
        Field Required
      </div>
    </div>
    <div class="form-group" :class="{ 'form-group--error': $v.bio.$error }">
      <label class="form__label">Bio</label>
      <textarea
        class="textarea has-fixed-size"
        v-model.trim="$v.bio.$model"
        placeholder="small Bio"
        rows="5"
      />
      <div class="error" v-if="!$v.bio.required && $v.bio.$dirty">
        Field Required
      </div>
      <div class="error" v-if="!$v.bio.max && $v.bio.$dirty">
        Bio should contain at most 200 characters
      </div>
    </div>
    <div class="form-group">
      <label class="form__label">Member Since</label>
      <input class="input" :value="member_since" type="text" readonly />
    </div>
    <div
      class="form-group"
      :class="{ 'form-group--error': $v.location.$error }"
    >
      <label class="form__label">Location</label>
      <input class="input" v-model.trim="$v.location.$model" type="text"  placeholder="location"/>
      <div class="error" v-if="!$v.location.required && $v.location.$dirty">
        Field Required
      </div>
    </div>
    <div class="form-group" :class="{ 'form-group--error': $v.URL.$error }">
      <label class="form__label">URL</label>
      <input class="input" v-model.trim="$v.URL.$model" type="text" placeholder="url"/>
      <div class="error" v-if="!$v.URL.required && $v.URL.$dirty">
        Field Required
      </div>
      <div class="error" v-if="!$v.URL.url && $v.URL.$dirty">
        Wrong Format
      </div>
    </div>
    <button class="button is-primary is-outlined" :disabled="$v.$invalid" @click="submit">
      Submit
    </button>
  </div>
</template>

<script>
import { required, maxLength, url} from "vuelidate/lib/validators"

export default {
  data() {
    return {
      name: "",
      bio: "",
      member_since: "",
      location: "",
      URL: "",
    }
  },
  methods: {
    submit: function() {
      if (this.$v.$invalid) return

      this.$http
        .put(this.$api_endpoints.USERS + this.$store.getters.getUsername, {
          name: this.name,
          bio: this.bio,
          location: this.location,
          url: this.URL,
        })
        .then(() => {
          this.$awn.success("Information Changed Successfully")
          this.trackSuccess()
          this.clear()
        })
        .catch(() => {
          this.$awn.alert("A problem occured, try again")
        })
    },
    trackSuccess() {
      this.$gtag.event("Account_information_change", {
        event_category: "Information_Change",
        event_label: "Info Changed",
        value: 1,
      })
    },
    clear() {
      this.name = ""
      this.bio = ""
      this.location = ""
      this.URL = ""
      this.$v.$reset()
    },
  },
  validations: {
    name: {
      required,
    },
    bio: {
      required,
      max: maxLength(255),
    },
    location: {
      required,
    },
    URL: {
      required,
      url,
    },
  },
}
</script>

<style lang="scss" scoped>
.page {
  margin-left: 5em;
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
}
</style>
