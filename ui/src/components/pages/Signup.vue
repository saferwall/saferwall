<template>
  <div class="container">
    <form class="form" novalidate="true" @submit.prevent="handleSubmit">
      <h1 class="signup">Create Your Account</h1>
      <div
        class="field entry"
        :class="{
          valid: !$v.username.$invalid,
          'not-valid': $v.username.$error,
        }"
      >
        <p class="control has-icons-left has-icons-right">
          <input
            v-focus
            required
            class="input"
            id="username"
            type="text"
            v-model.trim="$v.username.$model"
            placeholder="Username (e.g. John123)"
            autocomplete="username"
            @keyup.enter="handleSubmit"
          />
          <span class="icon is-small is-left">
            <i class="fas fa-user"></i>
          </span>
        </p>
        <div v-show="$v.username.$dirty">
          <span v-show="!$v.username.required" class="error"
            >Username is required</span
          >
        </div>
      </div>
      <div
        class="field entry"
        :class="{
          valid: !$v.email.$invalid,
          'not-valid': $v.email.$error,
        }"
      >
        <p class="control has-icons-left has-icons-right">
          <input
            required
            class="input"
            id="email"
            type="email"
            v-model.trim="$v.email.$model"
            placeholder="Email (name@example.com)"
            autocomplete="email"
            @keyup.enter="handleSubmit"
          />
          <span class="icon is-small is-left">
            <i class="fas fa-envelope"></i>
          </span>
        </p>
        <div v-show="$v.email.$dirty">
          <span v-show="!$v.email.required" class="error"
            >Email is required</span
          >

          <span v-show="!$v.email.email" class="error">Email is not valid</span>
          <!-- Add email in use constraint following backend implementation -->
        </div>
      </div>
      <div
        class="field entry"
        :class="{
          valid: !$v.password.$invalid,
          'not-valid': $v.password.$error,
        }"
      >
        <div>
          <p class="control has-icons-left has-icons-right">
            <input
              required
              class="input"
              id="password"
              :type="showPassword ? 'text' : 'password'"
              v-model.trim="$v.password.$model"
              placeholder="Password (Minimum 8 characters)"
              autocomplete="new-password"
              @keyup.enter="handleSubmit"
            />

            <button
              class="show-hide"
              @click.prevent="showPassword = !showPassword"
            >
              <svg
                v-if="showPassword"
                key="show-toggle"
                width="24"
                height="24"
                viewBox="0 0 24 24"
              >
                <path d="M0 0h24v24H0z" fill="none" />
                <path
                  d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"
                />
              </svg>
              <svg
                v-else
                key="show-toggle"
                width="24"
                height="24"
                viewBox="0 0 24 24"
              >
                <path
                  d="M0 0h24v24H0zm0 0h24v24H0zm0 0h24v24H0zm0 0h24v24H0z"
                  fill="none"
                />
                <path
                  d="M12 7c2.76 0 5 2.24 5 5 0 .65-.13 1.26-.36 1.83l2.92 2.92c1.51-1.26 2.7-2.89 3.43-4.75-1.73-4.39-6-7.5-11-7.5-1.4 0-2.74.25-3.98.7l2.16 2.16C10.74 7.13 11.35 7 12 7zM2 4.27l2.28 2.28.46.46C3.08 8.3 1.78 10.02 1 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l.42.42L19.73 22 21 20.73 3.27 3 2 4.27zM7.53 9.8l1.55 1.55c-.05.21-.08.43-.08.65 0 1.66 1.34 3 3 3 .22 0 .44-.03.65-.08l1.55 1.55c-.67.33-1.41.53-2.2.53-2.76 0-5-2.24-5-5 0-.79.2-1.53.53-2.2zm4.31-.78l3.15 3.15.02-.16c0-1.66-1.34-3-3-3l-.17.01z"
                />
              </svg>
            </button>
            <span class="icon is-small is-left">
              <i class="fas fa-lock"></i>
            </span>
          </p>
        </div>
        <div v-show="$v.password.$dirty">
          <span v-show="!$v.password.required" class="error"
            >Password is required</span
          >

          <span v-show="!$v.password.minLength" class="error"
            >Must be at least
            {{ $v.password.$params.minLength.min }} characters</span
          >
        </div>
      </div>
      <p :class="{ 'tos-required': $v.terms.$dirty && !$v.terms.sameAs }">
        <input v-model="$v.terms.$model" type="checkbox" id="tos" />
        <label for="tos"
          >&nbsp;I agree to the&nbsp;<router-link
            to="/tos"
            class="has-text-link"
            >Terms of Service</router-link
          ></label
        >
      </p>
      <button class="login" type="submit">
        Create Account
      </button>
    </form>
    <hr />
    <h3 class="already-member">
      Already have an account?
      <router-link :to="this.$routes.LOGIN.path" class="has-text-link"
        >Sign in</router-link
      >
    </h3>
  </div>
</template>

<script>
import {
  required,
  minLength,
  email,
  helpers,
  sameAs,
} from "vuelidate/lib/validators"

const usernameValid = helpers.regex("username", /^[a-zA-Z0-9]{1,20}$/)

export default {
  data() {
    return {
      username: "",
      email: "",
      password: "",
      terms: false,
      showPassword: false,
    }
  },
  methods: {
    handleSubmit() {
      this.$v.$touch()
      if (this.$v.$invalid) {
        this.$awn.alert("Please correct all highlighted errors and try again")
      } else {
        this.$http
          .post(this.$api_endpoints.USERS, {
            username: this.username,
            email: this.email,
            password: this.password,
          })
          .then((response) => {
            this.errored = false
            this.track()
            this.$router.push({
              path: this.$routes.LOGIN.path,
              query: {
                confirm: "email",
              },
            })
          })
          .catch(
            // server responded with a status code that falls out of the range of 2xx
            (error) => {
              this.$awn.alert(error.response.data.verbose_msg)
            },
          )
      }
    },
    track() {
      this.$gtag.event("sign_up", {
        method: "email",
      })
    },
  },
  validations: {
    username: {
      required,
      usernameValid,
    },
    email: {
      required,
      email,
    },
    password: {
      required,
      minLength: minLength(8),
    },
    terms: {
      sameAs: sameAs(() => true),
    },
  },
}
</script>

<style lang="scss" scoped>
@mixin bounce {
  display: inline-block;
  animation-duration: 0.3s;
  animation-iteration-count: 3;
  animation-name: bounce;
  animation-timing-function: ease-in-out;
}
.container {
  margin-bottom: 4em;
  background-color: white;
  border: 1px solid #33333330;
  border-radius: 0.25rem;
  padding-bottom: 2em;
}
.form {
  display: grid;
  align-items: center; /* align-self every label item vertically in its row!*/
  justify-content: center;
  width: 100%;
  padding: 2em;
  color: #333333;
  font-size: 16px;
}
hr {
  display: block;
  border: inherit;
  width: 22em;
  margin-top: 0.2em;
  margin-bottom: 1em;
}
.signup {
  font-size: 1em;
  font-weight: 600;
  margin-bottom: 1em;
  text-align: center;
}

.already-member {
  font-size: 0.9em;
  text-align: center;
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
.field > div > p {
  width: 340px;
  display: flex;
  flex-direction: row;
}

/* animate tos on error */
.tos-required {
  input {
    outline: 2px solid #bb0000;
  }
  label {
    @include bounce;
  }
}
@keyframes bounce {
  0% {
    transform: translateX(0);
  }
  50% {
    transform: translateX(4px);
  }
  100% {
    transform: translateX(0);
  }
}

.field:hover,
.field:focus-within {
  .show-hide,
  .entry {
    border-color: #3333336b;
    outline: none;
  }
}

.field > div > p > input {
  border-radius: 0.25rem 0 0 0.25rem !important;
  border-right: 0 !important;
}

.show-hide {
  width: 12%;
  background: none;
  border-radius: 0 0.25rem 0.25rem 0;
  border-width: 1px 1px 1px 0;
  border-style: solid;
  border-color: #33333335;
  box-shadow: inset -6px 2px 4px 0 hsla(0, 0%, 0%, 0.03);
}

.show-hide > svg {
  fill: #33333380;
}
.error {
  font-size: 0.8em;
  color: #bb0000;
  font-weight: 100;
}

.form .entry {
  min-height: 10px;
  color: #333333;
  background: none;
  padding: 0.5em;
  font-size: inherit;
  border-radius: 0.25rem;
  transition: border 0.1s ease;
}

.form .entry:focus {
  outline: none;
}
.login {
  cursor: pointer;
  border-radius: 0.25rem;
  font-size: inherit;
  padding: 0.7em;
  font-weight: 600;
  color: white;
  background-color: #18a096;
  border: none;
  margin-top: 10px;
}

.form label {
  font-weight: 500;
  font-size: 1rem;
}
.form label[for="tos"] {
  font-weight: normal;
  vertical-align: middle;
  font-size: 0.9em;
}
</style>
