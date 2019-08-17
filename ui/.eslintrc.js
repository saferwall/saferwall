module.exports = {
  extends: ["plugin:vue/essential", "standard", "plugin:prettier/recommended"],
  plugins: ["vue", "prettier", "standard"],
  rules: {
    // don't require .vue extension when importing
    "import/extensions": [
      "off",
      "always",
      {
        js: "never",
        vue: "never"
      }
    ],
    "no-new-wrappers": "off"
  }
};
