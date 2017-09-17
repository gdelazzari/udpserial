<template>
  <div class="configuration">
    <button class="uk-button uk-button-secondary" @click="onBtnApply()" style="float: right" :disabled="reloading == true">Apply configuration</button>
    <h2>Configured ports</h2>
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
      })
    },
    loadList() {
      this.$http.get('/api/ports').then(response => {
        this.ports = response.data
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
