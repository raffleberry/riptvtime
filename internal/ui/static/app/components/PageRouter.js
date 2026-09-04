import { RouterView } from "../vue.js"
import { Navigation } from "./Navigation.js"

const PageRouter = {
  components: {
    Navigation,
    RouterView,
  },
  props: {},
  setup() {
    return {}
  },

  template: /* HTML */ `
    <div class="vh-100 d-flex flex-column overflow-hidden">
      <Navigation></Navigation>

      <div class="flex-grow-1 d-flex flex-column mt-3 px-3 overflow-auto">
        <RouterView></RouterView>
      </div>
    </div>
  `,
}

export { PageRouter }
