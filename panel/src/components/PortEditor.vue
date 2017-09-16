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
                    <select class="uk-select uk-form-width-medium">
                      <option v-for="freePortName in freePortNames">
                        {{ freePortName }}
                      </option>
                    </select>
                  </div>
                </div>
                <div class="uk-margin">
                  <label class="uk-form-label" for="form-horizontal-text">Baudrate</label>
                  <div class="uk-form-controls">
                    <select class="uk-select uk-form-width-medium">
                      <option v-for="baudrate in baudrates">
                        {{ baudrate }}
                      </option>
                    </select>
                  </div>
                </div>
                <div class="uk-margin uk-grid-small uk-child-width-auto">
                  <label class="uk-form-label" for="form-horizontal-text">Data bits</label>
                  <div class="uk-form-controls">
                    <label><input class="uk-radio" type="radio" name="databitsRadio" checked> 8</label>
                    <label><input class="uk-radio" type="radio" name="databitsRadio"> 9</label>
                  </div>
                </div>
                <div class="uk-margin uk-grid-small uk-child-width-auto">
                  <label class="uk-form-label" for="form-horizontal-text">Stop bits</label>
                  <div class="uk-form-controls">
                    <label><input class="uk-radio" type="radio" name="stopbitsRadio"> 0</label>
                    <label><input class="uk-radio" type="radio" name="stopbitsRadio" checked> 1</label>
                  </div>
                </div>
                <div class="uk-margin">
                  <label class="uk-form-label" for="form-horizontal-text">Packet separator (empty for <i>auto</i>)</label>
                  <div class="uk-form-controls">
                    <input class="uk-input uk-form-width-small" type="text" placeholder="auto">
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
                    <select class="uk-select uk-form-width-medium">
                      <option v-for="ip in listenIPs">
                        {{ ip.ip }} ({{ ip.description }})
                      </option>
                    </select>
                  </div>
                </div>
                <div class="uk-margin">
                  <label class="uk-form-label" for="form-horizontal-text">UDP input (listening) port</label>
                  <div class="uk-form-controls">
                    <input class="uk-input uk-form-width-small" type="text" placeholder="5000">
                  </div>
                </div>
                <div class="uk-margin">
                  <label class="uk-form-label" for="form-horizontal-text">UDP output address</label>
                  <div class="uk-form-controls">
                    <input class="uk-input uk-form-width-small" type="text" placeholder="192.168.1.2">
                  </div>
                </div>
                <div class="uk-margin">
                  <label class="uk-form-label" for="form-horizontal-text">UDP output port</label>
                  <div class="uk-form-controls">
                    <input class="uk-input uk-form-width-small" type="text" placeholder="5000">
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'addport',
  data () {
    return {
      port: {},
      freePortNames: [],
      baudrates: [],
      listenIPs: []
    }
  },
  created () {
    this.$http.get('/api/freePortNames')
      .then(response => {
        this.freePortNames = response.data;
      })
    this.$http.get('/api/baudrates')
      .then(response => {
        this.baudrates = response.data;
      })
    this.$http.get('/api/listenIPs')
      .then(response => {
        this.listenIPs = response.data;
      })
  }
}
</script>

<style scoped>

</style>
