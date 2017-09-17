<script>
import { Line } from 'vue-chartjs'

export default Line.extend({
  name: 'portrate',
  data() {
    return {
      updater: null,
      colors: [
        'rgba(3, 201, 169, 0.7)',
        'rgba(68, 108, 179, 0.7)',
        'rgba(44, 62, 80, 0.7)',
        'rgba(231, 76, 60, 0.7)',
        'rgba(102, 51, 153, 0.7)'
      ]
    }
  },
  mounted() {
    this.initData()
  },
  destroyed() {
    clearInterval(this.updater)
    this._chart.destroy()
  },
  methods: {
    getRandomInt() {
      return Math.floor(Math.random() * (50))
    },
    initData() {
      this.$http.get('/api/statistics').then(response => {
        var data = {
          labels: [],
          datasets: []
        }

        var colorId = 0
        for (var portName in response.data.ports) {
          if (response.data.ports.hasOwnProperty(portName)) {
            data.datasets.push({
              label: portName,
              backgroundColor: this.colors[colorId],
              data: [],
              lineTension: 0.1,
              pointRadius: 2
            })
          }
          colorId++
          if (colorId >= this.colors.length) {
            colorId = 0
          }
        }

        var options = {
          responsive: true,
          maintainAspectRatio: false,
          animation: {
            easing: 'easeOutQuint'
          },
          scales: {
            yAxes: [{
              ticks: {
                min: 0,
                suggestedMax: 500
              }
            }]
          }
        }

        this.renderChart(data, options)
        this.appendStatistics(response.data)
        this.updater = setInterval(this.updateStatistics, 1000)
      })
    },
    updateStatistics() {
      if (this._chart.ctx == null) {
        clearInterval(this.updater)
        return
      }

      this.$http.get('/api/statistics').then(response => {
        this.appendStatistics(response.data)
      })
    },
    appendStatistics(data) {
      var willShift = false

      if (this._chart.data.labels.length >= 30) {
        willShift = true
      }

      for (var i = 0; i < this._chart.data.datasets.length; i++) {
        var dataset = this._chart.data.datasets[i]
        if (data.ports[dataset.label] != null) {
          var stats = data.ports[dataset.label]
          dataset.data.push(stats.udp2serialRate + stats.serial2udpRate)
          if (willShift) {
            dataset.data.shift()
          }
        }
      }

      if (willShift == false) {
        this._chart.data.labels.push('')
      }

      this._chart.update()
    }
  }
})
</script>

<style scoped>

</style>
