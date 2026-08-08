import { PageRouter } from "./components/PageRouter.js"
import { Confirm } from "./components/Confirm/Confirm.js"
import { PAGE, routes } from "./utils.js"
import {
  computed,
  createApp,
  createPinia,
  createRouter,
  createWebHistory,
  onMounted,
  ref,
  watch,
} from "./vue.js"
import { Notify } from "./components/Notify/Notify.js"

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  document.title = to.name || "Tv"
  next()
})

const app = createApp({
  components: {
    PageRouter,
    Confirm,
    Notify,
  },

  setup() {
    onMounted(() => {})

    return {}
  },
})

app.use(router)
app.use(createPinia())
app.mount("#app")
