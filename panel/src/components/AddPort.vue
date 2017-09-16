<template>
  <div class="addport">
    <h2>Add new port</h2>

    <porteditor></porteditor>

    <br><br>
    <button class="uk-button uk-button-primary">Confirm</button>
    <router-link class="uk-button uk-button-default" to="/configuration">Cancel</router-link>
  </div>
</template>

<script>
import PortEditor from '@/components/PortEditor'

export default {
  name: 'addport',
  components: {
    'porteditor': PortEditor
  },
  data () {
    return {
      ports: []
    }
  },
  created () {

  }
}
</script>

<style scoped>

</style>
