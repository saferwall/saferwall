<template>
  <div class="page">
    <div
      class="form-group"
      :class="{ 'form-group--error': $v.oldPassword.$error }"
    >
      <label class="form__label">Password</label>
      <input
        class="input"
        v-model.trim="$v.oldPassword.$model"
        type="password"
      />
      <div class="error" v-if="!$v.oldPassword.required && $v.oldPassword.$dirty">
        Old Password is required.
      </div>
    </div>
    <div
      class="form-group"
      :class="{ 'form-group--error': $v.password.$error }"
    >
      <label class="form__label">Password</label>
      <input class="input" v-model.trim="$v.password.$model" type="password" />
      <div class="error" v-if="!$v.password.required && $v.password.$dirty">
        Password is required.
      </div>
      <div class="error" v-if="!$v.password.minLength && $v.password.$dirty">
        Password must have at least
        {{ $v.password.$params.minLength.min }} letters.
      </div>
    </div>
    <div
      class="form-group"
      :class="{ 'form-group--error': $v.repeatPassword.$error }"
    >
      <label class="form__label">Repeat password</label>
      <input
        class="input"
        v-model.trim="$v.repeatPassword.$model"
        type="password"
      />
      <div class="error" v-if="!$v.repeatPassword.sameAsPassword && $v.repeatPassword.$dirty">
        Passwords must be identical.
      </div>
    </div>
    <button class="button is-primary is-outlined" :disabled="$v.$invalid">Submit</button>
  </div>
</template>

<script>
import { required, minLength, sameAs } from "vuelidate/lib/validators"

export default {
  data() {
    return {
      oldPassword: "",
      password: "",
      repeatPassword: "",
    }
  },
  validations: {
    oldPassword: {
      required,
    },
    password: {
      required,
      minLength: minLength(6),
    },
    repeatPassword: {
      sameAsPassword: sameAs("password"),
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

  .button{
      margin-top: 0.5em;
      width: 7em;
  }

  .form-group--error{
      .input{
          border-color: rgb(243, 54, 21);
      }
  }
}
</style>
