import { Navigation } from "./Navigation.js";
import { ref, RouterView } from "../vue.js";
import { handleSearch } from "../tabs/Search.js";

export const searchTerm = ref('')

const SearchBox = {
    components: {
    },
    props: {
    },

    setup() {
      return {
          searchTerm,
          handleSearch,
        }
    },

    template: `
    <div class="row">
      <div class="input-group search-container">
        <input v-model="searchTerm"
          type="text" class="form-control search-input"
          placeholder="Search..."
          @keyup.enter="handleSearch(searchTerm)">
        <button class="btn btn-outline-primary" type="button" id="searchButton"
          @click="handleSearch(searchTerm)">
          <i class="bi bi-search"></i>
        </button>
        
      </div>
    </div>
    `
}

export { SearchBox };
