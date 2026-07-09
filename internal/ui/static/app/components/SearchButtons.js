import { ref } from "../vue.js";
import { handleSearchNxtBtn, handleSearchPrvBtn, searchLoading, searchResults } from "../tabs/Search.js";


const SearchButtons = {
    components: {
    },
    props: {
    },
    setup() {
        return {
            searchLoading,
            searchResults,
            handleSearchNxtBtn,
            handleSearchPrvBtn,
        }
    },

    template: `
    <div class="d-flex flex-column justify-content-center align-items-center">
        <div class="input-group mt-3 justify-content-center">
            <button
                :disabled="searchResults.Page === 1 || searchLoading"
                type="button"
                class="btn btn-outline-primary" @click="handleSearchPrvBtn">
                <i class="bi bi-arrow-left"></i>
            </button>
            <div class="mx-3 d-flex flex-column align-items-center">
              <div>{{ searchResults.TotalResults }} Results</div>
              <div>{{ searchResults.Page }} / {{ searchResults.TotalPages }}</div>
            </div>
            <button
                :disabled="searchResults.Page === searchResults.TotalPages || searchLoading"
                type="button"
                class="btn btn-outline-primary" @click="handleSearchNxtBtn">
                <i class="bi bi-arrow-right"></i>
            </button>
        </div>
    </div>
    `
}

export { SearchButtons };
