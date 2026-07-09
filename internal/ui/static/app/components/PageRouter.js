import { Navigation } from "./Navigation.js";
import { RouterView } from "../vue.js";
import { currentPage, PAGE } from "../utils.js";
import { SearchBox } from "./SearchBox.js";
import { searchResults } from "../tabs/Search.js";
import { SearchButtons } from "./SearchButtons.js";


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
        return {
          PAGE,
          currentPage,

          searchResults
        }
    },

    template: `
    <div class="container-fluid vh-100 d-flex flex-column overflow-hidden">
      <Navigation></Navigation>

      <SearchBox class="my-2" v-if="currentPage.path === PAGE.SEARCH.path">
      </SearchBox>
      
      <div class="flex-grow-1 d-flex flex-row mt-3 overflow-auto">
        <RouterView></RouterView>
      </div>

      <SearchButtons class="my-2"
        v-if="currentPage.path === PAGE.SEARCH.path && searchResults.TotalResults > 0">
      </SearchButtons>
    </div>
`
}

export { PageRouter };
