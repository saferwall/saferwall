<template>
  <div class="page">
    <div class="notification">
      <div class="columns">
        <div class="column is-1 icn">
          <i class="icon is-large fas fa-exclamation-circle fa-3x"></i>
        </div>
        <div class="column">
          <h2>This Action is permanent</h2>
          <p>
            Deleting a user account is not a reversible change. Any information
            associated with this account will be removed, and cannot be
            recovered.
          </p>
        </div>
      </div>
    </div>
    <label class="checkbox">
      <input type="checkbox" v-model="valid" />
      I understand that this action cannot be undone.
    </label>
    <button class="button is-danger" :disabled="!valid">
      Delete Account
    </button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      valid: false,
    }
  },
}
</script>

<style lang="scss" scoped>
.page {
  display: grid;
  .notification {
    width: fit-content;
    margin-left: 5em;
    h2 {
      color: red;
      font-weight: 500;
      font-size: 1.5em;
    }
    .icon {
      color: red;
      margin: auto;
    }
    .icn {
      display: flex;
      align-items: center;
    }
  }
  .checkbox {
    margin: auto;
    margin-bottom: 1em;
    margin-left: 5em;
  }
  .button {
    width: fit-content;
    margin-left: 5em;
  }
}
</style>
