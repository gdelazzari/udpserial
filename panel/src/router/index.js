import Vue from 'vue'
import Router from 'vue-router'
import Statistics from '@/components/Statistics'
import Configuration from '@/components/Configuration'
import AddPort from '@/components/AddPort'
import EditPort from '@/components/EditPort'
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
      component: Statistics
    },
    {
      path: '/configuration',
      component: Configuration
    },
    {
      path: '/configuration/add',
      component: AddPort
    },
    {
      path: '/configuration/edit/:name',
      component: EditPort
    },
    {
      path: '/diagnostic',
      component: Diagnostic
    }
  ],
  linkActiveClass: 'uk-active'
})

// mode: 'history'
