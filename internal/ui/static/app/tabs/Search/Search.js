import { ENDPOINT, PAGE, theme } from "../../utils.js";
import { onMounted, ref, storeToRefs } from "../../vue.js";
import { SearchTile } from "./SearchTile.js";
import { useSeriesOpts } from "./optsStore.js";
import { useSearchStore } from "./searchStore.js";

const Search = {
    props: {
    },
    components: {
        SearchTile
    },
    setup: (props) => {
        
        const store = useSearchStore()

        const { loading, searchTerm, pageCur, results, resultsCnt } = storeToRefs(store)

        const { selected } = storeToRefs(useSeriesOpts())

        onMounted(() => {

        });

      
        return {
            loading,
            searchTerm,
            results,
            resultsCnt,
            pageCur,
        }
    },
    template: `
    <div class="d-flex flex-grow-1 flex-column">
        <div v-if="loading" class="d-flex justify-content-center align-items-center"
            style="min-height: 50vh;">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        <div v-else-if="resultsCnt === 0" class="d-flex justify-content-center align-items-center" style="min-height: 50vh;">
            <h2>Nothing</h2>
        </div>
        <div v-else class="col">
            <SearchTile class="mb-3" v-for="tv in results[pageCur]" :key="tv.Id" :tv="tv"></SearchTile>
        </div>
    </div>
    `
}
export { Search };

