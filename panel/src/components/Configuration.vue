<template>
  <div class="configuration">
    <button class="uk-button uk-button-secondary" @click="onBtnApply()" style="float: right" :disabled="reloading == true">Apply configuration</button>
    <h2>Configured ports</h2>
    <p class="uk-text-muted" v-if="ports.length == 0">No configured ports, start by adding a new one</p>
    <ul class="uk-list uk-list-large uk-list-striped">
      <li v-for="port in ports">
        {{ port }}
        <div class="list-buttons">
          <button class="uk-button uk-button-default uk-button-small" @click="onBtnEdit(port)">Edit</button>
          <button class="uk-button uk-button-danger uk-button-small" @click="onBtnUnlink(port)">Unlink</button>
        </div>
      </li>
    </ul>
    <router-link class="uk-button uk-button-primary" to="/configuration/add">Add new port</router-link>
  </div>
</template>

<script>
import UIkit from 'uikit'

export default {
  name: 'configuration',
  data () {
    return {
      ports: [],
      reloading: false
    }
  },
  created () {
    this.loadList()
  },
  methods: {
    onBtnUnlink(name) {
      this.$http.delete('/api/ports/' + name).then(response => {
        this.loadList()
      })
    },
    onBtnEdit(name) {
      this.$router.push("/configuration/edit/" + name)
    },
    onBtnApply() {
      this.reloading = true
      this.$http.get('/api/reloadConfigAndRestartThreads').then(response => {
        this.reloading = false
        UIkit.notification('Configuration reload completed', {pos: 'top-right', status: 'success'});
      })
    },
    loadList() {
      this.$http.get('/api/ports').then(response => {
        this.ports = response.data
        if (this.ports == null) {
          this.ports = []
        }
      })
    }
  }
}
</script>

<style scoped>
.list-buttons {
  float: right
}
</style>
