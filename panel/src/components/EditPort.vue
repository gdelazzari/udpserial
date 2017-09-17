<template>
  <div class="editport">
    <h2>Edit port configuration</h2>

    <porteditor :editing="true" :portName="$route.params.name" :onConfirm="onEditorConfirm" :onCancel="onEditorCancel"></porteditor>
  </div>
</template>

<script>
import PortEditor from '@/components/PortEditor'

export default {
  name: 'editport',
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
      this.$http.put('/api/ports/' + this.$route.params.name, portConfig).then(response => {
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
