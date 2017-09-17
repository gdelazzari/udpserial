<template>
  <div class="porteditor">
    <div class="uk-child-width-expand@s" uk-grid>
      <div>
        <div class="uk-card uk-card-small uk-card-default uk-card-body">
          <div class="uk-card-header">
            <h3 class="uk-card-title">Serial port settings</h3>
          </div>
          <div class="uk-card-body">
            <form class="uk-form-stacked">
              <div class="uk-margin">
                <label class="uk-form-label" for="form-horizontal-text">Port name</label>
                <div class="uk-form-controls">
                  <select class="uk-select uk-form-width-medium" v-model="port.name" :disabled="editing == true">
                    <option v-for="freePortName in freePortNames">
                      {{ freePortName }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="uk-margin">
                <label class="uk-form-label" for="form-horizontal-text">Baudrate</label>
                <div class="uk-form-controls">
                  <select class="uk-select uk-form-width-medium" type="number" v-model="port.baudrate">
                    <option v-for="baudrate in baudrates">
                      {{ baudrate }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="uk-margin uk-grid-small uk-child-width-auto">
                <label class="uk-form-label" for="form-horizontal-text">Data bits</label>
                <div class="uk-form-controls">
                  <select class="uk-select uk-form-width-small" type="number" v-model="port.databits">
                    <option>5</option>
                    <option>6</option>
                    <option>7</option>
                    <option>8</option>
                    <option>9</option>
                  </select>
                </div>
              </div>
              <div class="uk-margin uk-grid-small uk-child-width-auto">
                <label class="uk-form-label" for="form-horizontal-text">Stop bits</label>
                <div class="uk-form-controls">
                  <select class="uk-select uk-form-width-small" type="number" v-model="port.stopbits">
                    <option>1</option>
                    <option>2</option>
                  </select>
                </div>
              </div>
              <div class="uk-margin">
                <label class="uk-form-label" for="form-horizontal-text">Packet separator (empty for <i>auto</i>)</label>
                <div class="uk-form-controls">
                  <input class="uk-input uk-form-width-small" type="text" v-model="port.packetSeparator" placeholder="auto">
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
      <div>
        <div class="uk-card uk-card-small uk-card-default uk-card-body">
          <div class="uk-card-header">
            <h3 class="uk-card-title">UDP stream settings</h3>
          </div>
          <div class="uk-card-body">
            <form class="uk-form-stacked">
              <div class="uk-margin">
                <label class="uk-form-label" for="form-horizontal-text">UDP input (listening) address</label>
                <div class="uk-form-controls">
                  <select class="uk-select uk-form-width-medium" v-model="port.udpInputIP">
                    <option v-for="ip in listenIPs" :value="ip.ip">
                      {{ ip.ip }} ({{ ip.description }})
                    </option>
                  </select>
                </div>
              </div>
              <div class="uk-margin">
                <label class="uk-form-label" for="form-horizontal-text">UDP input (listening) port</label>
                <div class="uk-form-controls">
                  <input class="uk-input uk-form-width-small" type="number" placeholder="5000" v-model="port.udpInputPort">
                </div>
              </div>
              <div class="uk-margin">
                <label class="uk-form-label" for="form-horizontal-text">UDP output address</label>
                <div class="uk-form-controls">
                  <input class="uk-input uk-form-width-small" type="text" placeholder="x.x.x.x" v-model="port.udpOutputIP">
                </div>
              </div>
              <div class="uk-margin">
                <label class="uk-form-label" for="form-horizontal-text">UDP output port</label>
                <div class="uk-form-controls">
                  <input class="uk-input uk-form-width-small" type="number" placeholder="5000" v-model="port.udpOutputPort">
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    <br><br>
    <button class="uk-button uk-button-primary" @click="onBtnConfirm()">Confirm</button>
    <button class="uk-button uk-button-default" @click="onBtnCancel()">Cancel</button>
  </div>
</template>

<script>
export default {
  name: 'addport',
  data () {
    return {
      port: {
        name: "",
    		baudrate: 115200,
    		databits: 8,
    		stopbits: 1,
    		packetSeparator: "",
        udpInputIP: "0.0.0.0",
    		udpInputPort: 5000,
    		udpOutputIP: "localhost",
    		udpOutputPort: 5000
      },
      freePortNames: [],
      baudrates: [],
      listenIPs: []
    }
  },
  props: ['editing', 'onConfirm', 'onCancel', 'portName'],
  created () {
    this.$http.get('/api/baudrates').then(response => {
      this.baudrates = response.data
    })
    this.$http.get('/api/listenIPs').then(response => {
      this.listenIPs = response.data
    })
    if (this.editing == false) {
      this.$http.get('/api/freePortNames').then(response => {
        this.freePortNames = response.data
        this.port.name = response.data[0]
      })
    } else {
      this.freePortNames = [this.portName]
      this.$http.get('/api/ports/' + this.portName).then(response => {
        this.port = response.data
      })
    }
  },
  methods: {
    onBtnCancel(event) {
      this.onCancel()
    },
    onBtnConfirm(event) {
      this.port.baudrate = parseInt(this.port.baudrate)
      this.port.databits = parseInt(this.port.databits)
      this.port.stopbits = parseInt(this.port.stopbits)
      this.port.udpInputPort = parseInt(this.port.udpInputPort)
      this.port.udpOutputPort = parseInt(this.port.udpOutputPort)
      this.onConfirm(this.port)
    }
  }
}
</script>

<style scoped>

</style>
