<template>
  <div class="page">
    <div
      class="form-group"
      :class="{ 'form-group--error': $v.oldEmail.$error }"
    >
      <label class="form__label">Old Email</label>
      <input class="input" v-model.trim="$v.oldEmail.$model" type="text" />
      <div class="error" v-if="!$v.oldEmail.required && $v.oldEmail.$dirty">
        Field required.
      </div>
      <div class="error" v-if="!$v.oldEmail.email && $v.oldEmail.$dirty">
        Invalid email format.
      </div>
    </div>
    <div class="form-group" :class="{ 'form-group--error': $v.email.$error }">
      <label class="form__label">New Email</label>
      <input class="input" v-model.trim="$v.email.$model" type="text" />
      <div class="error" v-if="!$v.email.required && $v.email.$dirty">
        Field is required.
      </div>
      <div class="error" v-if="!$v.email.email && $v.email.$dirty">
        Invalid email format.
      </div>
      <div class="error" v-if="!$v.email.notSame && $v.email.$dirty">
        Same Email!
      </div>
    </div>
    <button class="button is-primary is-outlined" :disabled="$v.$invalid">
      Submit
    </button>
  </div>
</template>

<script>
import { required, email, sameAs, not } from "vuelidate/lib/validators"

export default {
  data() {
    return {
      oldEmail: "",
      email: "",
    }
  },
  validations: {
    oldEmail: {
      required,
      email,
    },
    email: {
      required,
      email,
      notSame : not(sameAs('oldEmail')),
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
