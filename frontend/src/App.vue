<script setup>
import Layout from './components/Layout.vue'
import FaucetForm from './components/FaucetForm.vue'
</script>

<template>
  <Layout>
    <FaucetForm />
  </Layout>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #141414;
}

* {
  box-sizing: border-box;
}
</style>
