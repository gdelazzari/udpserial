<template>
  <div id="app">
    <nav class="uk-navbar-container" uk-navbar>
      <div class="uk-navbar-left">
        <img class="uk-navbar-item uk-logo" src="./assets/logo.png"></img>
        <ul class="uk-navbar-nav">
          <router-link tag="li" to="/statistics"><a>Statistics</a></router-link>
          <router-link tag="li" to="/configuration"><a>Configuration</a></router-link>
          <router-link tag="li" to="/diagnostic"><a>Diagnostic</a></router-link>
        </ul>
      </div>
      <div class="uk-navbar-right">
        <div class="uk-navbar-item">
            <div>
              By Giacomo De Lazzari</br>
              <a href="mailto:gdelazzari@aol.com">gdelazzari@aol.com</a>
            </div>
        </div>
      </div>
    </nav>
    <div class="uk-container" id="main-container">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app'
}
</script>

<style>
#main-container {
  padding-top: 40px;
}
</style>
