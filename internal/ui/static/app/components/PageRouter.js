import { Navigation } from "./Navigation.js"
import { computed, RouterView, storeToRefs, useRoute } from "../vue.js"
import { PAGE } from "../utils.js"

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
