<template>
  <div class="diagnostic">
    <button class="uk-button uk-button-secondary" @click="onBtnRefresh()" style="float: right" :disabled="refreshing == true">Refresh</button>
    <h2>System log</h2>
    <div id="logDiv" class="uk-panel uk-panel-scrollable uk-height-large">
      <p id="logP" v-html="logText"></p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'diagnostic',
  data () {
    return {
      logText: '',
      refreshing: false
    }
  },
  created() {
    this.onBtnRefresh()
  },
  methods: {
    onBtnRefresh() {
      this.refreshing = true
      this.$http.get('/api/systemLog').then(response => {
        this.logText = response.data.replace(/(?:\r\n|\r|\n)/g, '<br />');
        setTimeout(function() {
          var logDiv = document.getElementById("logDiv")
          logDiv.scrollTop = logDiv.scrollHeight
        }, 10)
        this.refreshing = false
      })
    }
  }
}
</script>

<style scoped>
#logP {
  font-family: monospace;
  font-size: 14px;
}
</style>
