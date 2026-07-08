import { Navigation } from "./Navigation.js";
import { RouterView } from "../vue.js";


const PageRouter = {
    components: {
        Navigation,
        RouterView,
    },
    props: {
    },
    setup() {
        return {
        }
    },

    template: `
    <div class="container-fluid vh-100 d-flex flex-column overflow-hidden">
      <Navigation></Navigation>
      <div class="flex-grow-1 d-flex flex-row mt-3 overflow-auto">
        <RouterView></RouterView>
      </div>
    </div>
`
}

export { PageRouter };
