<template>
  <div class="configuration">
    <h2>Configured ports</h2>
    <ul class="uk-list uk-list-large uk-list-striped">
      <li v-for="port in ports">
        {{ port }}
        <div class="list-buttons">
          <button class="uk-button uk-button-default uk-button-small">Edit</button>
          <button class="uk-button uk-button-danger uk-button-small">Unlink</button>
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
      ports: []
    }
  },
  created () {
    this.$http.get('/api/ports')
      .then(response => {
        this.ports = response.data;
      })
  }
}
</script>

<style scoped>
.list-buttons {
  float: right
}
</style>
