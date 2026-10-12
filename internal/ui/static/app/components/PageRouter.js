import { RouterView } from "../vue.js"
import { isMobile } from "../utils.js"
import { Navigation } from "./Navigation.js"

const PageRouter = {
  components: {
    Navigation,
    RouterView,
  },
  props: {},
  setup() {
    return { isMobile }
  },

  template: /* HTML */ `
    <div class="vh-100 d-flex flex-column overflow-hidden">
      <Navigation></Navigation>

      <div
        class="flex-grow-1 d-flex flex-column mt-2 mt-md-3 px-2 px-md-3 overflow-auto"
        :class="{ 'rtt-has-bottom-nav': isMobile }"
      >
        <RouterView></RouterView>
      </div>
    </div>
  `,
}

export { PageRouter }
