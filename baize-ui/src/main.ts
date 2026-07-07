// @ts-nocheck
import { createApp } from 'vue'
import Antd from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import VxeUI from 'vxe-pc-ui'
import VXETable from 'vxe-table'
import 'ant-design-vue/dist/reset.css'
import 'vxe-pc-ui/lib/style.css'
import 'vxe-table/lib/style.css'

import '@/assets/styles/index.scss'

import App from './App.vue'
import store from './store/index'
import router from './router/index'
import directive from './directive/index'
import plugins from './plugins/index'
import { download } from '@/utils/request'
import 'virtual:svg-icons-register'
import SvgIcon from '@/components/SvgIcon/index.vue'
import elementIcons from '@/components/SvgIcon/svgicon'
import './permission'
import { useDict } from '@/utils/dict'
import { parseTime, resetForm, addDateRange, handleTree, selectDictLabel } from '@/utils/ruoyi'
import { handleProps } from '@/utils/baize'
import Pagination from '@/components/Pagination/index.vue'
import RightToolbar from '@/components/RightToolbar/index.vue'
import TreeSelect from '@/components/TreeSelect/index.vue'
import DictTag from '@/components/DictTag/index.vue'

const app = createApp(App)

app.config.globalProperties.useDict = useDict
app.config.globalProperties.download = download
app.config.globalProperties.parseTime = parseTime
app.config.globalProperties.resetForm = resetForm
app.config.globalProperties.handleTree = handleTree
app.config.globalProperties.handleProps = handleProps
app.config.globalProperties.addDateRange = addDateRange
app.config.globalProperties.selectDictLabel = selectDictLabel

app.component('DictTag', DictTag)
app.component('Pagination', Pagination)
app.component('TreeSelect', TreeSelect)
app.component('RightToolbar', RightToolbar)
app.component('svg-icon', SvgIcon)

app.use(router)
app.use(store)
app.use(Antd, { locale: zhCN })
app.use(plugins)
app.use(elementIcons)
app.use(VxeUI)
app.use(VXETable)

directive(app)

app.mount('#app')


