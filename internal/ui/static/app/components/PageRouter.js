import { Navigation } from "./Navigation.js";
import { computed, RouterView, storeToRefs, useRoute } from "../vue.js";
import { PAGE } from "../utils.js";
import { SearchBox } from "../tabs/Search/SearchBox.js";
import { SearchButtons } from "./SearchButtons.js";
import { useSearchStore } from "../tabs/Search/searchStore.js";


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

      <div class="flex-grow-1 d-flex flex-column mt-3 mx-2 overflow-auto">
        <RouterView></RouterView>
      </div>

    </div>
  `
}

export { PageRouter };
