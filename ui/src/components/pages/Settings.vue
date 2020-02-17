<template>
  <div class="tile is-vertical box">
    <p class="title">Account Settings:</p>
    <div class="columns">
      <div class="column is-one-quarter items">
        <div class="bt" :class="{ active: page === 0 }">
          <button class="button is-light" @click="page = 0">
            Account Information
          </button>
          <hr />
        </div>
        <div class="bt" :class="{ active: page === 1 }">
          <button class="button is-light" @click="page = 1">
            Change Email Address
          </button>
          <hr />
        </div>
        <div class="bt" :class="{ active: page === 2 }">
          <button class="button is-light" @click="page = 2">
            Change Password
          </button>
          <hr />
        </div>
        <div class="bt" :class="{ active: page === 3 }">
          <button class="button is-light" @click="page = 3">
            Delete Account
          </button>
          <hr />
        </div>
      </div>
      <div class="column">
        <ChangeAccountInfo v-if="page === 0" />
        <ChangeEmail v-if="page === 1" />
        <ChangePassword v-if="page === 2" />
        <DeleteAccount v-if="page === 3" />
      </div>
    </div>
  </div>
</template>

<script>
import ChangeAccountInfo from "../elements/settings/ChangeAccountInfo"
import ChangeEmail from "../elements/settings/ChangeEmail"
import ChangePassword from "../elements/settings/ChangePassword"
import DeleteAccount from "../elements/settings/DeleteAccount"

export default {
  components: {
    ChangeAccountInfo,
    ChangeEmail,
    ChangePassword,
    DeleteAccount,
  },
  data() {
    return {
      page: 0,
    }
  },
}
</script>

<style lang="scss" scoped>
.tile {
  padding-bottom: 2em;
  padding-right: 4em;
  padding-left: 4em;

  .title {
    margin-bottom: 20px;
  }
  .columns {
    margin-left: 20px;

    .items {
      display: grid;
      width: max-content;
      height: max-content;
      text-align: left;

      .bt {
        display: inline-flex;
        height: 35px;
        margin-bottom: 0;
        margin-top: 0;

        .button {
          background-color: transparent;
          width: 200px;
          justify-content: left;
        }
        hr {
          background-color: grey;
          border: none;
          display: block;
          height: 35px;
          width: 2px;
          margin: 0;
        }
      }

      .active {
        button {
          color: #00d1b2;
          font-weight: 600;
        }
        hr {
          background-color: #00d1b2;
        }
      }

      .bt:hover {
        button {
          color: #00d1b2;
          font-weight: 600;
        }
        hr {
          background-color: #00d1b2;
        }
      }
    }
  }
}
</style>
