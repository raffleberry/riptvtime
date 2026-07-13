import { Navigation } from "./Navigation.js";
import { computed, RouterView, storeToRefs, useRoute } from "../vue.js";
import { PAGE } from "../utils.js";
import { SearchBox } from "./SearchBox.js";
import { SearchButtons } from "./SearchButtons.js";
import { useSearchStore } from "../stores/search.js";


const PageRouter = {
    components: {
        Navigation,
        RouterView,
        SearchBox,
        SearchButtons,
    },
    props: {
    },
    setup() {
        const store = useSearchStore()
        const { resultsCnt } = storeToRefs(store)

        const r = useRoute()

        const curPath = computed(() => r.path)


        return {
          PAGE,
          curPath,

          resultsCnt
        }
    },

    template: `
    <div class="container-fluid vh-100 d-flex flex-column overflow-hidden">
      <Navigation></Navigation>

      <SearchBox class="my-2" v-if="curPath === PAGE.SEARCH.path">
      </SearchBox>
      
      <div class="flex-grow-1 d-flex flex-row mt-3 overflow-auto">
        <RouterView></RouterView>
      </div>

      <SearchButtons class="my-2"
        v-if="curPath === PAGE.SEARCH.path && resultsCnt > 0">
      </SearchButtons>
    </div>
`
}

export { PageRouter };
