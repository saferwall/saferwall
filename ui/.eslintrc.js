module.exports = {
  extends: ["plugin:vue/essential", "standard"],
  plugins: ["vue"],
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
