<template>
  <div class="statistics">
    <h2>Live transmission rates</h2>
    <portrategraph></portrategraph>
  </div>
</template>

<script>
import PortRateGraph from '@/components/PortRateGraph'

export default {
  name: 'statistics',
  components: {
    'portrategraph': PortRateGraph
  },
  data () {
    return {

    }
  },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped>

</style>
