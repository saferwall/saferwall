import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/pages/Home'
import Upload from '@/components/pages/Upload'
import Scanning from '@/components/pages/Scanning'
import Antivirus from '@/components/pages/Antivirus'
import Summary from '@/components/pages/Summary'
import Strings from '@/components/pages/Strings'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home,
      meta: {title: 'Home'}
    },
    {
      path: '/upload',
      name: 'Upload',
      component: Upload,
      meta: {title: 'Upload'}
    },
    {
      path: '/scanning',
      name: 'Scanning',
      component: Scanning,
      meta: {title: 'Scanning'}
    },
    {
      path: '/antivirus/:hash',
      name: 'Antivirus',
      component: Antivirus,
      meta: {title: 'Antivirus'}
    },
    {
      path: '/summary/:hash',
      name: 'Summary',
      component: Summary,
      meta: {title: 'Summary'}
    },
    {
      path: '/strings/:hash',
      name: 'Strings',
      component: Strings,
      meta: {title: 'Strings'}
    }
  ]

})
