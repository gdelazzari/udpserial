import Vue from 'vue'
import Router from 'vue-router'
import Statistics from '@/components/Statistics'
import Configuration from '@/components/Configuration'
import AddPort from '@/components/AddPort'
import Diagnostic from '@/components/Diagnostic'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      redirect: '/statistics'
    },
    {
      path: '/statistics',
      name: 'Statistics',
      component: Statistics
    },
    {
      path: '/configuration',
      name: 'Configuration',
      component: Configuration
    },
    {
      path: '/configuration/add',
      name: 'Add port',
      component: AddPort
    },
    {
      path: '/diagnostic',
      name: 'Diagnostic',
      component: Diagnostic
    }
  ],
  linkActiveClass: 'uk-active'
})

// mode: 'history'
