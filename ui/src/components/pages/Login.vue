<template>
  <div class="container">
    <form
      novalidate="true"
      class="form"
      @submit.prevent="handleSubmit"
      v-if="emailConfirmed"
    >
      <h1 class="signin">Sign In</h1>
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
            v-model.trim="username"
            placeholder="Username"
            autocomplete="username"
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
          valid: !$v.password.$invalid,
          'not-valid': $v.password.$error,
        }"
      >
        <p class="control has-icons-left">
          <input
            required
            class="input"
            id="password"
            type="password"
            v-model.trim="password"
            placeholder="Password"
            autocomplete="current-password"
          />
          <span class="icon is-small is-left">
            <i class="fas fa-lock"></i>
          </span>
        </p>
        <div v-show="$v.password.$dirty">
          <span v-show="!$v.password.required" class="error"
            >Password is required</span
          >
        </div>
      </div>
      <button class="login" type="submit">Sign In</button>
      <h3 class="forgot">
        <router-link :to="this.$routes.FORGOT_PWD.path" class="has-text-link"
          >Forgot password?</router-link
        >
        or
        <span class="confirm" @click="emailConfirmed = false"
          >didn't confirm email?</span
        >
      </h3>
    </form>
    <confirm v-if="!emailConfirmed" @sent="emailSent" />
    <hr />
    <h3 class="not-member" v-if="emailConfirmed">
      Not a member?
      <router-link :to="this.$routes.SIGNUP.path">Sign up</router-link>
    </h3>
  </div>
</template>

<script>
import { required, helpers } from "vuelidate/lib/validators"
import ConfirmEmailForm from "../elements/ConfirmEmailForm"
import { mapActions } from "vuex"

const usernameValid = helpers.regex("username", /^[a-zA-Z0-9]{1,20}$/)

export default {
  data() {
    return {
      username: "",
      password: "",
      confirmationWarning:
        "Confirm your email by clicking the verification link we just sent to your inbox.",
      emailConfirmed: true,
    }
  },
  components: {
    confirm: ConfirmEmailForm,
  },
  mounted() {
    if (this.$route.query.confirm) {
      this.$awn.warning(this.confirmationWarning)
    }
  },
  methods: {
    ...mapActions(["updateLoggedIn", "updateUsername"]),
    handleSubmit() {
      this.$v.$touch()
      if (this.$v.$invalid) {
        this.$awn.alert("Please correct all highlighted errors and try again")
      } else {
        this.$http
          .post(this.$api_endpoints.AUTH_LOGIN, {
            username: this.username,
            password: this.password,
          })
          .then((response) => {
            // We store a second cookie which contains the payload only.
            // The cookie which contains the auth token is stored on a httpOnly cookie.
            this.$cookies.set("JWTPayload", response.data.token.split(".")[1])
            this.updateLoggedIn(response.data.token.split(".")[1])
            this.updateUsername(response.data.token.split(".")[1])

            this.track()

            if (this.$route.params.nextUrl != null) {
              this.$router.push(this.$route.params.nextUrl)
            } else {
              this.$router.push(this.$routes.HOME.path)
            }
          })
          .catch((e) => {
            this.$awn.alert(e.response.data.verbose_msg)
            if (
              e.response.data.verbose_msg ===
              "Account not confirmed, please confirm your email !"
            ) {
              this.emailConfirmed = false
            }
          })
      }
    },
    emailSent() {
      this.emailConfirmed = true
    },
    track() {
      this.$gtag.event("login", {
        method: "email",
      })
    },
  },
  validations: {
    username: {
      required,
      usernameValid,
    },
    password: {
      required,
    },
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
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

.container {
  margin-bottom: 4em;
  background-color: white;
  border-radius: 0.25rem;
  border: 1px solid #33333330;
  padding-bottom: 2em;
}
.form {
  display: grid;
  line-height: 1em;
  align-items: center; /* align-self every label item vertically in its row!*/
  justify-content: center;
  width: 100%;
  grid-row-gap: 0.1em;
  padding: 1em 2em;
  color: #333333;
  font-size: 16px;
  
}
hr{
display: block;
    border: inherit;
    width: 22em;
    margin-top: 0.5em;
    margin-bottom: 1em;
}

.input-container {
  display: flex;
  height: 80px;
  flex-direction: column;
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
  font-size: 0.8em;
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

.forgot,
.signin {
  padding-top: 14px;
  margin: auto;
}

.signin {
  font-size: 1em;
  font-weight: 600;
  margin-bottom: 1em;
  text-align: center;
}

.not-member {
  font-size: 0.9em;
  text-align: center;
}

.form .entry {
  width: 340px;
  min-height: 45px;
  color: #333333;
  background: none;
  padding: 0.5em;
  font-size: 0.8em;
  border-radius: 0.25rem;
  transition: border 0.1s ease;
}

.form .entry:focus {
  outline: none;
}

.login {
  cursor: pointer;
  font-size: 0.8em;
  border-radius: 0.25rem;
  padding: 1em;
  font-weight: 600;
  color: white;
  background-color: #18a096;
  border: none;
  margin-top: 0.8em;
}

.form label {
  font-weight: 600;
  font-size: 0.8em;
}
.has-text-link {
  font-size: 0.8em;
}
.confirm {
  font-size: 0.8em;
  cursor: pointer;
  font-weight: 400;
  color: #3273dc;
}
</style>
