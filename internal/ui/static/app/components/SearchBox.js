import { Navigation } from "./Navigation.js";
import { ref, RouterView } from "../vue.js";
import { useSearchStore } from "../stores/search.js";

const SearchBox = {
    components: {
    },
    props: {
    },

    setup() {

      const searchTerm = ref('')

      const store = useSearchStore()
      const { onSearch } = store

      return {
          searchTerm,
          onSearch,
        }
    },

    template: `
    <div class="row">
      <div class="input-group search-container">
        <input v-model="searchTerm"
          type="text" class="form-control search-input"
          placeholder="Search..."
          @keyup.enter="onSearch(searchTerm)">
        <button class="btn btn-outline-primary" type="button" id="searchButton"
          @click="onSearch(searchTerm)">
          <i class="bi bi-search"></i>
        </button>
        
      </div>
    </div>
    `
}

export { SearchBox };
