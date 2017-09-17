<script>
import { Line } from 'vue-chartjs'

export default Line.extend({
  name: 'portrate',
  data() {
    return {
      datacollection: {
        labels: ['', ''],
        datasets: [
          {
            label: 'TTL-1',
            backgroundColor: 'rgba(231, 76, 60, 0.7)',
            data: [this.getRandomInt(), this.getRandomInt()],
            lineTension: 0.1
          },
          {
            label: 'TTL-2',
            backgroundColor: 'rgba(102, 51, 153, 0.7)',
            data: [this.getRandomInt(), this.getRandomInt()],
            lineTension: 0.1
          }
        ]
      },
      updater: null
    }
  },
  mounted() {
    this.renderChart(this.datacollection, {responsive: true, maintainAspectRatio: false, animation: {easing: 'easeOutQuint'}})
    this.updater = setInterval(this.appendValue, 1000)
  },
  destroyed() {
    clearInterval(this.updater)
    this._chart.destroy()
  },
  methods: {
    getRandomInt() {
      return Math.floor(Math.random() * (50))
    },
    appendValue() {
      if (this._chart.ctx == null) {
        clearInterval(this.updater)
        return
      }

      this._chart.data.datasets[0].data.push(this.getRandomInt())
      this._chart.data.datasets[1].data.push(this.getRandomInt())
      this._chart.data.labels.push('')

      this._chart.update()
    }
  }
})
</script>

<style scoped>

</style>
