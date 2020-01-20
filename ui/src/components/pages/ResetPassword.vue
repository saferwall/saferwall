<!-- refactor input/login markup/logic into reusable vue components -->

<template>
  <div class="container">
    <form class="form" novalidate="true" @submit.prevent="handleSubmit">
      <h1 class="heading">Reset Password</h1>
      <p class="instruction">
        Password must be at least 8 characters long.
      </p>
      <div
        class="field entry"
        :class="{
          valid: !$v.password.$invalid,
          'not-valid': $v.password.$error,
        }"
      >
        <div>
          <p class="control has-icons-left">
            <input
              v-focus
              required
              class="input"
              id="password"
              :type="showPassword ? 'text' : 'password'"
              v-model.trim="$v.password.$model"
              placeholder="New Password"
              autocomplete="new-password"
              @keyup.enter="handleSubmit"
            /><button
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
      <div
        class="field entry"
        :class="{
          valid: !$v.password.$invalid && !$v.repeatPassword.$invalid,
          'not-valid': $v.repeatPassword.$error,
        }"
      >
        <p class="control has-icons-left">
          <input
            required
            class="input"
            id="repeatPassword"
            type="password"
            v-model.trim="$v.repeatPassword.$model"
            placeholder="Retype New Password"
            autocomplete="new-password"
            @keyup.enter="handleSubmit"
          />
          <span class="icon is-small is-left">
            <i class="fas fa-lock"></i>
          </span>
        </p>
        <div v-show="$v.repeatPassword.$dirty">
          <span v-show="!$v.repeatPassword.sameAsPassword" class="error"
            >Passwords must be identical</span
          >
        </div>
      </div>
      <button class="reset-btn" type="submit">
        Reset Password
      </button>
    </form>
  </div>
</template>

<script>
import { required, minLength, sameAs } from "vuelidate/lib/validators"
import queryString from "query-string"

export default {
  data() {
    return {
      password: "",
      repeatPassword: "",
      showPassword: false,
      successMessage: "Password reset successful!",
    }
  },
  methods: {
    handleSubmit() {
      this.$v.$touch()
      if (this.$v.$invalid) {
        this.$awn.alert("Please correct all highlighted errors and try again")
      } else {
        const { token } = this.$route.query
        const postData = {
          password: this.password,
        }
        const axiosConfig = {
          headers: {
            "Content-Type": "application/json",
          },
          params: {
            token,
          },
          paramsSerializer: (params) =>
            queryString.stringify(params, { arrayFormat: "bracket" }),
        }
        this.$http
          .post(this.$api_endpoints.AUTH_CHANGE_PWD, postData, axiosConfig)
          .then((response) => {
            this.$awn.success(this.successMessage)
            this.$router.push(this.$routes.HOME.path)
          })
          .catch((error) => {
            this.$awn.alert(
              error.response.data.verbose_msg ||
                "An error occurred. Please try again later!",
            )
          })
      }
    },
  },
  validations: {
    password: {
      required,
      minLength: minLength(8),
    },
    repeatPassword: {
      sameAsPassword: sameAs("password"),
    },
  },
}
</script>

<style lang="scss" scoped>
.form {
  display: grid;
  text-align: center;
  align-items: center; /* align-self every label item vertically in its row!*/
  justify-content: center;
  width: min-content;
  padding: 2em;
  color: #333333;
  background-color: white;
  font-size: 16px;
  border-radius: 0.25rem;
  border: 1px solid #33333330;
}
.field {
  display: flex;
  height: 100%;
  flex-direction: column;
}


.field > div > p {
  width: 340px;
  display: flex;
  flex-direction: row;
}
.field:hover,
.field:focus-within {
  .show-hide,
  .entry {
    border-color: #3333336b;
    outline: none;
  }
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
.field > div > p > input {
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
  min-height: 10px;
  color: #333333;
  background: none;
  padding: 0.5em;
  font-size: inherit;
  border-radius: 0.25rem;
  transition: border 0.1s ease;
}
.error {
  font-size: 0.8em;
  color: #bb0000;
  font-weight: 100;
}

.heading {
  font-size: 1em;
  font-weight: 600;
  margin-bottom: 1em;
  text-align: center;
}

.instruction {
  font-size: 1em;
  line-height: 1.4em;
  font-weight: 300;
  margin-bottom: 10px;
}
.reset-btn {
  cursor: pointer;
  border-radius: 0.25rem;
  font-size: inherit;
  padding: 0.7em;
  font-weight: 600;
  color: white;
  background-color: #18a096;
  border: none;
}
</style>
