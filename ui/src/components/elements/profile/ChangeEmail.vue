<template>
  <div class="page">
    <div
      class="form-group"
      :class="{ 'form-group--error': $v.email.$error }"
    >
      <label class="form__label">Email</label>
      <input class="input" v-model.trim="$v.email.$model" type="text" />
      <div class="error" v-if="!$v.email.required && $v.email.$dirty">
        Field required.
      </div>
      <div class="error" v-if="!$v.email.email && $v.email.$dirty">
        Invalid email format.
      </div>
    </div>

    <div
      class="form-group"
      :class="{ 'form-group--error': $v.password.$error }"
    >
      <label class="form__label">Password</label>
      <input
        class="input"
        v-model.trim="$v.password.$model"
        type="password"
      />
      <div
        class="error"
        v-if="!$v.password.required && $v.password.$dirty"
      >
        Old Password is required.
      </div>
      <div
        class="error"
        v-if="!$v.password.minLength && $v.password.$dirty"
      >
        Password must have at least
        {{ $v.password.$params.minLength.min }} letters.
      </div>
    </div>

    <button
      class="button is-primary is-outlined"
      :disabled="$v.$invalid"
      @click="submit"
    >
      Submit
    </button>
  </div>
</template>

<script>
import { required, email, minLength } from "vuelidate/lib/validators"

export default {
  data() {
    return {
      email: "",
      password: "",
    }
  },
  methods: {
    submit: function() {
      if (this.$v.$invalid) return

      this.$http
        .put(
          this.$api_endpoints.USERS +
            this.$store.getters.getUsername +
            "/email",
          {
            email: this.email,
            password: this.password,
          },
        )
        .then(() => {
          this.$awn.success("Email Changed Successfully")
          this.trackSuccess()
          this.clear()
        })
        .catch(() => {
          this.$awn.alert("A problem occured, try again")
        })
    },
    trackSuccess() {
      this.$gtag.event("Email_change", {
        event_category: "Information_Change",
        event_label: "Email Changed",
        value: 1,
      })
    },
    clear() {
      this.email = ""
      this.password = ""
      this.$v.$reset()
    },
  },
  validations: {
    email: {
      required,
      email,
    },
    password: {
      required,
      minLength: minLength(8),
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
