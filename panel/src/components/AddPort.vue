<template>
  <div class="addport">
    <h2>Add new port configuration</h2>

    <porteditor :editing="false" :onConfirm="onEditorConfirm" :onCancel="onEditorCancel"></porteditor>
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

  },
  methods: {
    onEditorCancel() {
      this.$router.push("/configuration")
    },
    onEditorConfirm(portConfig) {
      this.$http.post('/api/ports', portConfig).then(response => {
        if (response.data != null) {
          this.$router.push("/configuration")
        }
      });
    }
  }
}
</script>

<style scoped>

</style>
