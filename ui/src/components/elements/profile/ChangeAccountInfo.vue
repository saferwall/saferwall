<template>
  <div class="page">
    <div class="form-group" :class="{ 'form-group--error': $v.name.$error }">
      <label class="form__label">Name</label>
      <input class="input" v-model.trim="$v.name.$model" type="text" />
      <div class="error" v-if="!$v.name.required && $v.name.$dirty">
        Field Required
      </div>
      <div class="error" v-if="!$v.name.min && $v.name.$dirty">
        Name should contain at least 5 characters
      </div>
    </div>
    <div class="form-group">
      <label class="form__label">Email</label>
      <input class="input" :value="email" type="text" readonly />
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
      <div class="error" v-if="!$v.bio.min && $v.bio.$dirty">
        Bio should contain at least 10 characters
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
      <input class="input" v-model.trim="$v.location.$model" type="text" />
      <div class="error" v-if="!$v.location.required && $v.location.$dirty">
        Field Required
      </div>
    </div>
    <div class="form-group" :class="{ 'form-group--error': $v.URL.$error }">
      <label class="form__label">URL</label>
      <input class="input" v-model.trim="$v.URL.$model" type="text" />
      <div class="error" v-if="!$v.URL.required && $v.URL.$dirty">
        Field Required
      </div>
      <div class="error" v-if="!$v.URL.url && $v.URL.$dirty">
        Wrong Format
      </div>
    </div>
    <button class="button is-primary is-outlined" :disabled="$v.$invalid">
      Submit
    </button>
  </div>
</template>

<script>
import { required, maxLength, minLength, url } from "vuelidate/lib/validators"

export default {
  data() {
    return {
      name: "walid",
      email: "walid@gmail.com",
      bio:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
      member_since: "12/10/2019",
      location: "France",
      URL: "/walidO",
    }
  },
  validations: {
    name: {
      required,
      min: minLength(5),
    },
    bio: {
      required,
      max: maxLength(200),
      min: minLength(10),
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
