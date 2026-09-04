import { Confirm } from "./components/Confirm/Confirm.js"
import { Notify } from "./components/Notify/Notify.js"
import { PageRouter } from "./components/PageRouter.js"
import { routes } from "./utils.js"
import { createApp, createPinia, createRouter, createWebHistory, onMounted } from "./vue.js"

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  if (to.name) {
    document.title = `${to.name} | riptvtime`
  } else {
    document.title = "riptvtime"
  }
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
