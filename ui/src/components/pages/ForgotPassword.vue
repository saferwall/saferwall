<template>
  <div class="container">
    <transition name="slide" mode="out-in">
      <notification
        type="is-success"
        @closeNotif="closeMessage()"
        v-if="succeeded"
      >
        {{ successMessage }}
      </notification>
    </transition>
    <transition name="slide" mode="out-in">
      <notification type="is-danger" @closeNotif="close()" v-if="errored">
        {{ errorMessage }}
      </notification>
    </transition>
    <form class="form" novalidate="true" @submit.prevent="handleSubmit">
      <h1 class="heading">Forgot Password?</h1>
      <p class="instruction">
        Enter your account email address and we will send you a link to reset
        your password.
      </p>
      <div
        class="input-container"
        :class="{
          valid: !$v.email.$invalid,
          'not-valid': $v.email.$error,
        }"
      >
        <label for="email">Email</label>

        <input
          v-focus
          required
          class="entry"
          id="email"
          type="email"
          v-model.trim="$v.email.$model"
          placeholder="name@example.com"
          autocomplete="email"
          @keyup.enter="handleSubmit"
        />
        <div v-show="$v.email.$dirty">
          <span v-show="!$v.email.required" class="error"
            >Email is required</span
          >
        </div>
      </div>
      <button class="reset-btn" type="submit">
        Request Password Reset
      </button>
    </form>
  </div>
</template>

<script>
import { required, email } from "vuelidate/lib/validators"
import Notification from "@/components/elements/Notification"
export default {
  data() {
    return {
      email: "",
      succeeded: false,
      errored: false,
      errorMessage: "",
      successMessage:
        "We've sent a password reset link to the email you specified",
    }
  },

  components: {
    notification: Notification,
  },
  methods: {
    handleSubmit() {
      this.$v.$touch()
      if (this.$v.$invalid) {
        this.errored = true
        this.errorMessage = "Please enter a valid email address"
      } else {
        this.$http
          .delete(this.api_endpoints.AUTH_CHANGE_PWD, {
            data: {
              email: this.email,
            },
          })
          .then((response) => {
            this.succeeded = true
            this.errored = false
          })
          .catch((error) => {
            this.errored = true
            this.succeeded = false
            this.errorMessage =
              error.response.data.verbose_msg ||
              "An error occurred. Please try again later!"
          })
      }
    },
    close() {
      this.errored = false
    },
    closeMessage() {
      this.succeeded = false
    },
  },
  validations: {
    email: {
      required,
      email,
    },
  },
}
</script>

<style lang="scss" scoped>
.form {
  display: grid;
  grid-template-rows: 1fr 1fr minmax(70px, min-content) 1fr;
  text-align: center;
  grid-row-gap: 1.5em;
  line-height: 2em;
  align-items: center; /* align-self every label item vertically in its row!*/
  justify-content: center;
  width: min-content;
  padding: 4em;
  color: #333333;
  background-color: white;
  font-size: 16px;
  border-radius: 0.25rem;
  border: 1px solid #33333330;
}
.input-container {
  display: flex;
  height: 100%;
  flex-direction: column;
}

.input-container > label {
  text-align: left;
}

.input-container > div {
  width: 340px;
  display: flex;
  flex-direction: row;
}
.input-container:hover,
.input-container:focus-within {
  .show-hide,
  .entry {
    border-color: #3333336b;
    outline: none;
  }
}

.input-container > div > input {
  border-radius: 0.25rem 0 0 0.25rem !important;
  border-right: 0 !important;
}

.valid {
  input {
    background: url("../../assets/imgs/check.svg") !important;
    background-repeat: no-repeat !important;
    background-position: 97% center !important;
    border: 1px solid #02a678 !important;
  }

  .show-hide {
    border-color: #02a678 !important;
  }
}

.not-valid {
  input {
    background: url("../../assets/imgs/error.svg") !important;
    background-repeat: no-repeat !important;
    background-position: 97% center !important;
    border: 1px solid #bb0000 !important;
  }
  .show-hide {
    border-color: #bb0000 !important;
  }
}
.form .entry:focus {
  outline: none;
}
.form .entry {
  width: 340px;
  min-height: 45px;
  color: #333333;
  background: none;
  border: 1px solid #33333335;
  padding: 0.5em;
  font-size: inherit;
  border-radius: 0.25rem;
  box-shadow: inset 6px 2px 4px 0 hsla(0, 0%, 0%, 0.03);
  transition: border 0.1s ease;
}
.error {
  font-size: 0.8em;
  color: #bb0000;
  font-weight: 100;
}

.heading {
  font-size: 1.5em;
}

.instruction {
  font-size: 1em;
  line-height: 1.4em;
  font-weight: 300;
}
.reset-btn {
  cursor: pointer;
  border-radius: 0.25rem;
  font-size: inherit;
  padding: 0.7em;
  font-weight: 600;
  color: white;
  background-color: #e7501d;
  border: none;
}
</style>
